package handlers

import (
	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Auth         *AuthHandler
	Users        *UsersHandler
	Roles        *RolesHandler
	Applications *ApplicationsHandler
	Dashboard    *DashboardHandler
	Audit        *AuditHandler
}

func New(db *pgxpool.Pool, cfg *config.Config) *Handlers {
	repos := repository.New(db)

	authSvc := services.NewAuthService(repos.Users, repos.Tokens, cfg)
	userSvc := services.NewUserService(repos.Users, repos.Roles)
	appSvc := services.NewApplicationService(repos.Applications, repos.Audit)

	return &Handlers{
		Auth:         &AuthHandler{svc: authSvc},
		Users:        &UsersHandler{svc: userSvc, audit: repos.Audit},
		Roles:        &RolesHandler{roles: repos.Roles, audit: repos.Audit},
		Applications: &ApplicationsHandler{svc: appSvc, audit: repos.Audit},
		Dashboard:    &DashboardHandler{users: repos.Users, apps: repos.Applications},
		Audit:        &AuditHandler{repo: repos.Audit},
	}
}
