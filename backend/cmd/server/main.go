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
	// Load .env in development; ignore error in production where env vars are set externally
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

	h := handlers.New(db, cfg)

	rl := middleware.NewRateLimiter(cfg.RateLimitReqs, cfg.RateLimitWindow)
	authRL := middleware.NewRateLimiter(10, cfg.RateLimitWindow) // strict: 10 req/window for auth

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger)
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
		// ── Public auth ──────────────────────────────────────────────────────────
		r.Route("/auth", func(r chi.Router) {
			r.Use(authRL.Middleware)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.RefreshToken)
			r.Post("/forgot-password", h.Auth.ForgotPassword)
			r.Post("/reset-password", h.Auth.ResetPassword)
		})

		// ── Authenticated routes ─────────────────────────────────────────────────
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(cfg.JWTSecret))

			// Self-service auth
			r.Post("/auth/logout", h.Auth.Logout)
			r.Get("/auth/me", h.Auth.Me)
			r.Post("/auth/change-password", h.Auth.ChangePassword)

			// ── Applicant-facing application routes ─────────────────────────
			r.Route("/applications", func(r chi.Router) {
				// Draft management
				r.Post("/draft", h.Applications.SaveDraft)
				r.Get("/draft/{formType}", h.Applications.GetDraft)
				r.Delete("/draft/{formType}", h.Applications.DeleteDraft)

				// User's own applications
				r.Get("/", h.Applications.List)
				r.Get("/{id}", h.Applications.Get)

				// Wizard flow
				r.Post("/{id}/signed-form", h.Applications.UploadSignedForm)
				r.Get("/{id}/signed-form", h.Applications.GetSignedForm)
				r.Post("/{id}/payment", h.Applications.UploadPaymentProof)
				r.Get("/{id}/payment", h.Applications.GetPayment)
				r.Post("/{id}/submit", h.Applications.Submit)
				r.Post("/{id}/respond", h.Applications.Respond)

				// Attachments
				r.Post("/{id}/attachments", h.Applications.UploadAttachment)
				r.Get("/{id}/attachments", h.Applications.ListAttachments)
				r.Delete("/{id}/attachments/{attachmentID}", h.Applications.DeleteAttachment)

				// PDF generation (future)
				r.Get("/{id}/pdf", h.Applications.GeneratePDF)
			})

			// Transactions (own)
			r.Get("/transactions", h.Applications.ListTransactions)
			r.Get("/transactions/{id}", h.Applications.GetTransaction)

			// ── Admin: dashboard + users + audit ────────────────────────────
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
			})

			// ── Admin: application review (admin + general_secretary) ────────
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

			// ── Super admin: roles + permissions ────────────────────────────
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
