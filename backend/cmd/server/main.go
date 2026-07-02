package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/database"
	"github.com/atenimedia-llc/ncs-online/backend/internal/handlers"
	appLogging "github.com/atenimedia-llc/ncs-online/backend/internal/logging"
	"github.com/atenimedia-llc/ncs-online/backend/internal/maintenance"
	appMetrics "github.com/atenimedia-llc/ncs-online/backend/internal/metrics"
	"github.com/atenimedia-llc/ncs-online/backend/internal/middleware"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	logCloser, err := appLogging.Configure(cfg.LogLevel, cfg.LogDir)
	if err != nil {
		log.Fatalf("logging: %v", err)
	}
	defer logCloser.Close()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	h, repos := handlers.New(db, cfg)
	// Start the GitHub-release sentinel so the admin "Smart Updates" page
	// can show update availability without an explicit user click.
	h.UpdatesSvc.StartSentinel(context.Background())
	// Auto-clear an expired maintenance window so the UI flips back to
	// "operational" without the operator returning to the page.
	maintenance.StartScheduleSentinel(context.Background(), h.SystemState, repos.Operator)
	locationClient := middleware.NewLocationClient(cfg.LocationServiceURL, cfg.LocationServiceToken, cfg.LocationTimeout)
	middleware.ConfigureLocationService(locationClient, cfg.JWTSecret, cfg.GeoAllowedCountries, cfg.GeoFailClosed)
	auditWriter := middleware.NewAuditWriter(repos.Audit, 4096)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = auditWriter.Close(ctx)
	}()

	rl := middleware.NewRateLimiter(cfg.RateLimitReqs, cfg.RateLimitWindow)
	authRL := middleware.NewRateLimiter(10, cfg.RateLimitWindow)

	r := chi.NewRouter()
	metricRegistry := appMetrics.New()

	r.Use(chimiddleware.RequestID)
	r.Use(metricRegistry.Middleware)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RejectAmbiguousPaths)
	r.Use(middleware.LimitRequestBody(25 << 20))
	r.Use(middleware.Logger)
	r.Use(middleware.AuditLogger(auditWriter))
	r.Use(middleware.HoneypotScanner)
	r.Use(middleware.MaintenanceMode(h.SystemState, cfg.MaintenanceEnabled, cfg.JWTSecret))
	r.Use(rl.Middleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(chimiddleware.Recoverer)

	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
	r.Get("/health", healthHandler)
	r.Get("/healthz", healthHandler)
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			response.Err(w, http.StatusServiceUnavailable, "NOT_READY", "Database is unavailable")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})
	r.Handle("/metrics", metricRegistry)
	r.Post("/ncs-ussd", h.USSD.ServeHTTP)
	r.Get("/api/v1/system/maintenance-status", func(w http.ResponseWriter, r *http.Request) {
		s := h.SystemState.Get()
		response.JSON(w, http.StatusOK, map[string]any{
			"public_cms":      s.Scoped(maintenance.ScopePublicCMS),
			"admin_dashboard": s.Scoped(maintenance.ScopeAdminDashboard),
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// ── Public auth (geo-blocked: Uganda only, no VPN) ────────────
		r.Route("/auth", func(r chi.Router) {
			r.Use(middleware.GeoBlocker)
			r.Use(authRL.Middleware)
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.RefreshToken)
			r.Post("/forgot-password", h.Auth.ForgotPassword)
			r.Post("/reset-password", h.Auth.ResetPassword)
		})

		// ── Public CMS (read-only, published content) ─────────────────
		r.Route("/cms", func(r chi.Router) {
			r.Get("/posts", h.CMS.ListPosts)
			r.Get("/posts/{slug}", h.CMS.GetPost)
			r.Get("/posts/{slug}/comments", h.CMS.ListPostComments)
			r.Get("/blog/categories", h.CMS.ListBlogCategories)
			r.Get("/events", h.CMS.ListEvents)
			r.Get("/events/{slug}", h.CMS.GetEvent)
			r.Get("/careers", h.CMS.ListCareers)
			r.Get("/careers/{id}", h.CMS.GetCareer)
			r.Get("/slides", h.CMS.ListSlides)
			r.Get("/slideshows/{slug}", h.CMS.GetSlideshow)
			r.Get("/menus/{name}", h.CMS.GetMenu)
			r.Get("/settings/{key}", h.CMS.GetSetting)
			r.Get("/fun-facts", h.CMS.ListFunFacts)
			r.Get("/faqs", h.CMS.ListFAQs)
			r.Get("/resources", h.CMS.ListResources)
			r.Get("/facilities", h.CMS.ListFacilities)
			r.Get("/associations", h.CMS.ListAssociations)
			r.Get("/invest", h.CMS.ListInvest)
			r.Get("/team", h.CMS.ListTeam)
			r.Get("/departments", h.CMS.ListInstitutionalDepartments)
			r.Post("/newsletter/subscribe", h.CMS.SubscribeNewsletter)
		})

		// ── Authenticated routes (geo-blocked: Uganda only, no VPN) ─────
		r.Group(func(r chi.Router) {
			r.Use(middleware.GeoBlocker)
			r.Use(middleware.Authenticate(cfg.JWTSecret))
			r.Use(middleware.ValidateGlobalSession(h.SystemState))
			r.Use(middleware.ValidateAuthenticatedUser(repos.Users))
			r.Use(middleware.EnforceIPAllowlist(repos.Security))

			// Self-service auth
			r.Post("/auth/logout", h.Auth.Logout)
			r.Get("/auth/me", h.Auth.Me)
			r.Post("/auth/change-password", h.Auth.ChangePassword)

			// PIN management
			r.Post("/auth/pin/set", h.Auth.SetPIN)
			r.Put("/auth/pin/change", h.Auth.ChangePIN)
			r.Post("/auth/pin/verify", h.Auth.VerifyPIN)
			r.Get("/me/contexts", h.Organisations.Contexts)
			r.Post("/organisation-invitations/accept", h.Organisations.AcceptInvitation)
			r.Get("/organisations/{organisationID}", h.Organisations.Get)
			r.Get("/me/sessions", h.Operator.Sessions)
			r.Post("/me/sessions/revoke", h.Operator.RevokeSession)

			// ── Applicant application routes ─────────────────────────
			r.Route("/applications", func(r chi.Router) {
				r.Post("/draft", h.Applications.SaveDraft)
				r.Get("/draft/{formType}", h.Applications.GetDraft)
				r.Delete("/draft/{formType}", h.Applications.DeleteDraft)

				r.Get("/", h.Applications.List)
				r.Get("/{id}", h.Applications.Get)

				r.Post("/{id}/signed-form", h.Applications.UploadSignedForm)
				r.Get("/{id}/signed-form", h.Applications.GetSignedForm)
				r.Post("/{id}/payment", h.Applications.UploadPaymentProof)
				r.Get("/{id}/payment", h.Applications.GetPayment)
				r.Post("/{id}/submit", h.Applications.Submit)
				r.Post("/{id}/respond", h.Applications.Respond)

				r.Post("/{id}/attachments", h.Applications.UploadAttachment)
				r.Get("/{id}/attachments", h.Applications.ListAttachments)
				r.Delete("/{id}/attachments/{attachmentID}", h.Applications.DeleteAttachment)

				r.Get("/{id}/pdf", h.Applications.GeneratePDF)
			})

			r.Get("/transactions", h.Applications.ListTransactions)
			r.Get("/transactions/{id}", h.Applications.GetTransaction)

			// ── Departments (read-only directory) ─────────────────────
			r.Get("/departments", h.Forms.ListDepartments)

			// ── Dynamic application portal (applicant) ────────────────
			r.Route("/portal/forms", func(r chi.Router) {
				r.Get("/open", h.Forms.PortalListOpen)
				r.Get("/{slug}", h.Forms.PortalGetForm)
				r.Post("/{templateID}/draft", h.Forms.PortalSaveDraft)
			})
			r.Route("/portal/submissions", func(r chi.Router) {
				r.Get("/{id}", h.Forms.PortalGetSubmission)
				r.Post("/{id}/submit", h.Forms.PortalSubmit)
				r.Post("/{id}/payment-proof", h.Forms.PortalUploadPaymentProof)
			})

			// Generic authenticated file upload (used by the dynamic form
			// dropzone for both applicants and admins).
			r.Post("/media/upload", h.CMS.UploadMedia)
			r.With(authRL.Middleware).Post("/cms/posts/{slug}/comments", h.CMS.SubmitPostComment)

			// ── Self-service activity log + security settings ─────────
			r.Get("/me/activities", h.Security.MyActivities)
			r.Route("/me/security", func(r chi.Router) {
				r.Get("/", h.Security.GetMySecurity)
				r.Post("/ip-whitelist", h.Security.AddIPWhitelist)
				r.Delete("/ip-whitelist/{id}", h.Security.RemoveIPWhitelist)
				r.Post("/2fa/enroll", h.Security.Enroll2FA)
				r.Post("/2fa/verify", h.Security.Verify2FA)
				r.Post("/2fa/disable", h.Security.Disable2FA)
			})

			r.Route("/nsmis", func(r chi.Router) {
				r.Get("/federations", h.NSMIS.ListFederations)
				r.Post("/federations", h.NSMIS.CreateFederation)
				r.Put("/federations/{federationID}", h.NSMIS.UpdateFederation)
				r.Get("/reporting-periods", h.NSMIS.ListPeriods)
				r.Post("/reporting-periods", h.NSMIS.CreatePeriod)
				r.Post("/reporting-periods/{periodID}/obligations", h.NSMIS.GenerateObligations)
				r.Get("/report-obligations", h.NSMIS.ListObligations)
				r.Put("/report-obligations/{obligationID}/governance-draft", h.NSMIS.SaveGovernanceDraft)
				r.Post("/reports/{reportID}/transitions", h.NSMIS.TransitionReport)
				r.Get("/dashboards/governance", h.NSMIS.GovernanceDashboard)
				r.Get("/dashboards/athletes", h.NSMIS.AthleteDashboard)
				r.Get("/dashboards/performance", h.NSMIS.PerformanceDashboard)
				r.Get("/dashboards/finance", h.NSMIS.FinanceDashboard)
				r.Get("/dashboards/talent", h.NSMIS.TalentDashboard)
				r.Get("/federations/{federationID}/documents", h.NSMIS.ListDocuments)
				r.Post("/federations/{federationID}/documents", h.NSMIS.UploadDocument)
				r.Get("/documents/{documentID}/download", h.NSMIS.DownloadDocument)
				r.Get("/jobs", h.NSMIS.ListJobs)
				r.Post("/jobs", h.NSMIS.QueueJob)
				r.Post("/jobs/{jobID}/retry", h.NSMIS.RetryJob)
				r.Route("/{resource:federation-officers|federation-memberships|athletes|competitions|medals|coaches|technical-officials|talent|scholarships|safeguarding-aggregates|safeguarding-cases|disbursements|accountabilities|equipment}", func(r chi.Router) {
					r.Get("/", h.NSMIS.ListDomain)
					r.Post("/", h.NSMIS.CreateDomain)
					r.Put("/{id}", h.NSMIS.UpdateDomain)
					r.Delete("/{id}", h.NSMIS.DeleteDomain)
				})
			})

			// ── Admin: dashboard + users + audit ─────────────────────
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles("super_admin", "admin"))

				r.Get("/admin/dashboard", h.Dashboard.Stats)
				r.Get("/admin/system/status", h.Operator.Status)
				r.Get("/admin/system/resources", h.Operator.Resources)
				r.Get("/admin/system/resources/ws", h.Operator.ResourceStream)
				r.Get("/admin/system/service-logs", h.Operator.ServiceLogs)

				r.Route("/admin/users", func(r chi.Router) {
					r.Get("/", h.Users.List)
					r.Post("/", h.Users.Create)
					r.Get("/{id}", h.Users.Get)
					r.Put("/{id}", h.Users.Update)
					r.Delete("/{id}", h.Users.Delete)
					r.Post("/{id}/activate", h.Users.Activate)
					r.Post("/{id}/deactivate", h.Users.Deactivate)
					r.Post("/{id}/reset-password", h.Users.ResetPassword)
					r.Post("/{id}/account-action", h.Users.AccountAction)
					r.Post("/{id}/roles", h.Users.AssignRole)
					r.Delete("/{id}/roles/{roleID}", h.Users.RemoveRole)
				})

				r.Route("/admin/audit-logs", func(r chi.Router) {
					r.Get("/", h.Audit.List)
					r.Get("/{id}", h.Audit.Get)
					r.Get("/export", h.Audit.Export)
				})

				// Admin needs read access to roles list (for assignment UI)
				r.Get("/admin/roles-list", h.Roles.List)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles("super_admin"))
				r.Put("/admin/system/maintenance", h.Operator.SetMaintenance)
				r.Post("/admin/system/actions", h.Operator.MaintenanceAction)
				r.Post("/admin/system/services/action", h.Operator.ServiceAction)
				r.Post("/admin/system/cache/flush", h.Operator.FlushCache)
				r.Post("/admin/system/sessions/revoke-all", h.Operator.RevokeAll)
				r.Get("/admin/storage-settings", h.CMS.GetStorageSettings)
				r.Put("/admin/storage-settings", h.CMS.UpdateStorageSettings)
				r.Post("/admin/storage-settings/google-drive/connect", h.CMS.ConnectGoogleDrive)
				r.Post("/admin/storage-settings/google-drive/exchange", h.CMS.ExchangeGoogleDriveCode)
				r.Post("/admin/storage-settings/google-drive/disconnect", h.CMS.DisconnectGoogleDrive)
				r.Get("/admin/system/backups", h.Backups.List)
				r.Get("/admin/system/backups/schema", h.Backups.DownloadSchema)
				r.Post("/admin/system/backups/schema/import", h.Backups.ImportSchema)
				r.Delete("/admin/system/backups/schema", h.Backups.DeleteSchema)
				r.Get("/admin/system/backups/jobs/export", h.Backups.ExportJobs)
				r.Delete("/admin/system/backups/jobs", h.Backups.ClearJobs)
				r.Delete("/admin/system/backups/jobs/{jobID}", h.Backups.DeleteJob)
				r.Get("/admin/system/backups/{id}/download", h.Backups.Download)
				r.Delete("/admin/system/backups/{id}", h.Backups.Delete)
				r.Post("/admin/system/backups/jobs", h.Backups.Queue)

				// Smart updates (GitHub release sentinel + in-place deploy).
				// Deploy/rollback require the Docker socket to be mounted into
				// the backend container — see Docs/runbooks/smart-updates.md.
				r.Get("/admin/system/updates", h.Updates.Status)
				r.Get("/admin/system/updates/settings", h.Updates.Settings)
				r.Put("/admin/system/updates/settings", h.Updates.SaveSettings)
				r.Post("/admin/system/updates/check", h.Updates.Check)
				r.Get("/admin/system/updates/deploy", h.Updates.DeployStatus)
				r.Post("/admin/system/updates/deploy", h.Updates.Deploy)
				r.Post("/admin/system/updates/rollback", h.Updates.Rollback)
			})

			// ── Admin: application review ─────────────────────────────
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles("super_admin", "admin", "general_secretary"))

				r.Route("/admin/applications", func(r chi.Router) {
					r.Get("/", h.Applications.AdminList)
					r.Get("/{id}", h.Applications.AdminGet)
					r.Post("/{id}/approve", h.Applications.Approve)
					r.Post("/{id}/reject", h.Applications.Reject)
					r.Post("/{id}/request-info", h.Applications.RequestInfo)
					r.Post("/{id}/verify-payment", h.Applications.VerifyPayment)
					r.Post("/{id}/reject-payment", h.Applications.RejectPayment)
				})

				r.Get("/admin/transactions", h.Applications.ListTransactions)

				// ── Admin: dynamic form templates (department-scoped) ──
				r.Route("/admin/forms", func(r chi.Router) {
					r.Get("/", h.Forms.AdminListTemplates)
					r.Post("/", h.Forms.AdminCreateTemplate)
					r.Get("/submissions", h.Forms.AdminListSubmissions)
					r.Get("/submissions/{id}", h.Forms.AdminGetSubmission)
					r.Post("/submissions/{id}/review", h.Forms.AdminReviewSubmission)
					r.Post("/submissions/{id}/verify-payment", h.Forms.AdminVerifySubmissionPayment)
					r.Get("/{id}", h.Forms.AdminGetTemplate)
					r.Put("/{id}", h.Forms.AdminUpdateTemplate)
					r.Delete("/{id}", h.Forms.AdminDeleteTemplate)
				})
			})

			// ── Admin: CMS management ─────────────────────────────────
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles("super_admin", "admin", "content_manager"))

				r.Route("/admin/cms/posts", func(r chi.Router) {
					r.Get("/", h.CMS.ListPosts)
					r.Post("/", h.CMS.CreatePost)
					r.Put("/{id}", h.CMS.UpdatePost)
					r.Delete("/{id}", h.CMS.DeletePost)
				})

				r.Route("/admin/cms/blog/categories", func(r chi.Router) {
					r.Get("/", h.CMS.ListBlogCategories)
					r.Post("/", h.CMS.CreateBlogCategory)
					r.Put("/{id}", h.CMS.UpdateBlogCategory)
					r.Delete("/{id}", h.CMS.DeleteBlogCategory)
				})

				r.Route("/admin/cms/comments", func(r chi.Router) {
					r.Get("/", h.CMS.ListCommentsModeration)
					r.Post("/{id}/approve", h.CMS.ApproveComment)
					r.Post("/{id}/flag", h.CMS.FlagComment)
					r.Delete("/{id}", h.CMS.DeleteComment)
				})

				r.Get("/admin/cms/newsletter/subscribers", h.CMS.ListNewsletterSubscribers)
				r.Get("/admin/cms/newsletter/subscribers/export", h.CMS.ExportNewsletterSubscribers)

				r.Route("/admin/cms/events", func(r chi.Router) {
					r.Get("/", h.CMS.ListEvents)
					r.Post("/", h.CMS.CreateEvent)
					r.Put("/{id}", h.CMS.UpdateEvent)
					r.Delete("/{id}", h.CMS.DeleteEvent)
				})

				r.Route("/admin/cms/careers", func(r chi.Router) {
					r.Get("/", h.CMS.ListCareers)
					r.Post("/", h.CMS.CreateCareer)
					r.Put("/{id}", h.CMS.UpdateCareer)
					r.Delete("/{id}", h.CMS.DeleteCareer)
				})

				r.Route("/admin/cms/slides", func(r chi.Router) {
					r.Get("/", h.CMS.ListSlides)
					r.Post("/", h.CMS.CreateSlide)
					r.Put("/{id}", h.CMS.UpdateSlide)
					r.Delete("/{id}", h.CMS.DeleteSlide)
				})
				r.Get("/admin/cms/slideshows/{slug}", h.CMS.GetSlideshow)
				r.Put("/admin/cms/slideshows/{slug}", h.CMS.UpdateSlideshow)

				r.Put("/admin/cms/menus/{name}", h.CMS.UpdateMenu)
				r.Put("/admin/cms/settings/{key}", h.CMS.UpdateSetting)

				r.Route("/admin/cms/fun-facts", func(r chi.Router) {
					r.Get("/", h.CMS.ListFunFacts)
					r.Post("/", h.CMS.CreateFunFact)
					r.Put("/{id}", h.CMS.UpdateFunFact)
					r.Delete("/{id}", h.CMS.DeleteFunFact)
				})

				r.Route("/admin/cms/faqs", func(r chi.Router) {
					r.Get("/", h.CMS.ListFAQs)
					r.Post("/", h.CMS.CreateFAQ)
					r.Put("/{id}", h.CMS.UpdateFAQ)
					r.Delete("/{id}", h.CMS.DeleteFAQ)
				})

				r.Route("/admin/cms/resources", func(r chi.Router) {
					r.Get("/", h.CMS.ListResources)
					r.Post("/", h.CMS.CreateResource)
					r.Put("/{id}", h.CMS.UpdateResource)
					r.Delete("/{id}", h.CMS.DeleteResource)
				})

				r.Route("/admin/cms/facilities", func(r chi.Router) {
					r.Get("/", h.CMS.ListFacilities)
					r.Post("/", h.CMS.CreateFacility)
					r.Put("/{id}", h.CMS.UpdateFacility)
					r.Delete("/{id}", h.CMS.DeleteFacility)
				})

				r.Route("/admin/cms/associations", func(r chi.Router) {
					r.Get("/", h.CMS.ListAssociations)
					r.Post("/", h.CMS.CreateAssociation)
					r.Put("/{id}", h.CMS.UpdateAssociation)
					r.Delete("/{id}", h.CMS.DeleteAssociation)
				})

				r.Route("/admin/cms/invest", func(r chi.Router) {
					r.Get("/", h.CMS.ListInvest)
					r.Post("/", h.CMS.CreateInvest)
					r.Put("/{id}", h.CMS.UpdateInvest)
					r.Delete("/{id}", h.CMS.DeleteInvest)
				})

				r.Route("/admin/cms/team", func(r chi.Router) {
					r.Get("/", h.CMS.ListTeam)
					r.Post("/", h.CMS.CreateTeamMember)
					r.Put("/{id}", h.CMS.UpdateTeamMember)
					r.Delete("/{id}", h.CMS.DeleteTeamMember)
				})

				r.Route("/admin/cms/departments", func(r chi.Router) {
					r.Get("/", h.CMS.ListInstitutionalDepartments)
					r.Post("/", h.CMS.CreateInstitutionalDepartment)
					r.Put("/{id}", h.CMS.UpdateInstitutionalDepartment)
					r.Delete("/{id}", h.CMS.DeleteInstitutionalDepartment)
				})

				r.Post("/admin/media/upload", h.CMS.UploadMedia)
			})

			// ── Super admin: roles + permissions ─────────────────────
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles("super_admin"))

				r.Route("/admin/roles", func(r chi.Router) {
					r.Get("/", h.Roles.List)
					r.Post("/", h.Roles.Create)
					r.Get("/{id}", h.Roles.Get)
					r.Put("/{id}", h.Roles.Update)
					r.Delete("/{id}", h.Roles.Delete)
					r.Post("/{id}/permissions", h.Roles.AssignPermission)
					r.Delete("/{id}/permissions/{permID}", h.Roles.RemovePermission)
				})

				r.Get("/admin/permissions", h.Roles.ListPermissions)
			})
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.Err(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "HTTP method is not allowed for this endpoint")
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server listening", "port", cfg.Port, "environment", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-quit
	slog.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	slog.Info("server stopped")
}
