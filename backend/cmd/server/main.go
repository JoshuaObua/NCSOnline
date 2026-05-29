package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/database"
	"github.com/atenimedia-llc/ncs-online/backend/internal/handlers"
	"github.com/atenimedia-llc/ncs-online/backend/internal/middleware"
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

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	h, repos := handlers.New(db, cfg)

	rl := middleware.NewRateLimiter(cfg.RateLimitReqs, cfg.RateLimitWindow)
	authRL := middleware.NewRateLimiter(10, cfg.RateLimitWindow)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.AuditLogger(repos.Audit))
	r.Use(rl.Middleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// ── Public auth (geo-blocked: Uganda only, no VPN) ────────────
		r.Route("/auth", func(r chi.Router) {
			r.Use(middleware.GeoBlocker)
			r.Use(authRL.Middleware)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.RefreshToken)
			r.Post("/forgot-password", h.Auth.ForgotPassword)
			r.Post("/reset-password", h.Auth.ResetPassword)
		})

		// ── Public CMS (read-only, published content) ─────────────────
		r.Route("/cms", func(r chi.Router) {
			r.Get("/posts", h.CMS.ListPosts)
			r.Get("/posts/{slug}", h.CMS.GetPost)
			r.Get("/events", h.CMS.ListEvents)
			r.Get("/events/{slug}", h.CMS.GetEvent)
			r.Get("/careers", h.CMS.ListCareers)
			r.Get("/careers/{id}", h.CMS.GetCareer)
			r.Get("/slides", h.CMS.ListSlides)
			r.Get("/menus/{name}", h.CMS.GetMenu)
			r.Get("/fun-facts", h.CMS.ListFunFacts)
			r.Get("/faqs", h.CMS.ListFAQs)
			r.Get("/resources", h.CMS.ListResources)
			r.Get("/facilities", h.CMS.ListFacilities)
			r.Get("/associations", h.CMS.ListAssociations)
			r.Get("/invest", h.CMS.ListInvest)
		})

		// ── Authenticated routes (geo-blocked: Uganda only, no VPN) ─────
		r.Group(func(r chi.Router) {
			r.Use(middleware.GeoBlocker)
			r.Use(middleware.Authenticate(cfg.JWTSecret))

			// Self-service auth
			r.Post("/auth/logout", h.Auth.Logout)
			r.Get("/auth/me", h.Auth.Me)
			r.Post("/auth/change-password", h.Auth.ChangePassword)

			// PIN management
			r.Post("/auth/pin/set", h.Auth.SetPIN)
			r.Put("/auth/pin/change", h.Auth.ChangePIN)
			r.Post("/auth/pin/verify", h.Auth.VerifyPIN)

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

			// ── Admin: dashboard + users + audit ─────────────────────
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRoles("super_admin", "admin"))

				r.Get("/admin/dashboard", h.Dashboard.Stats)

				r.Route("/admin/users", func(r chi.Router) {
					r.Get("/", h.Users.List)
					r.Post("/", h.Users.Create)
					r.Get("/{id}", h.Users.Get)
					r.Put("/{id}", h.Users.Update)
					r.Delete("/{id}", h.Users.Delete)
					r.Post("/{id}/activate", h.Users.Activate)
					r.Post("/{id}/deactivate", h.Users.Deactivate)
					r.Post("/{id}/roles", h.Users.AssignRole)
					r.Delete("/{id}/roles/{roleID}", h.Users.RemoveRole)
				})

				r.Route("/admin/audit-logs", func(r chi.Router) {
					r.Get("/", h.Audit.List)
					r.Get("/{id}", h.Audit.Get)
				})

				// Admin needs read access to roles list (for assignment UI)
				r.Get("/admin/roles-list", h.Roles.List)
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

				r.Put("/admin/cms/menus/{name}", h.CMS.UpdateMenu)

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

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening on :%s (env=%s)", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped")
}