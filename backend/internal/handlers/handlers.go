package handlers

import (
	"context"
	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/maintenance"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Auth          *AuthHandler
	Users         *UsersHandler
	Roles         *RolesHandler
	Applications  *ApplicationsHandler
	Dashboard     *DashboardHandler
	Audit         *AuditHandler
	CMS           *CMSHandler
	NSMIS         *NSMISHandler
	USSD          *USSDHandler
	Organisations *OrganisationsHandler
	Operator      *OperatorHandler
	SystemState   *maintenance.State
	Backups       *BackupsHandler
	Forms         *FormsHandler
	Security      *SecurityHandler
	Updates       *UpdatesHandler
	UpdatesSvc    *services.UpdatesService
	Analytics     *AnalyticsHandler
}

// New constructs all handlers and returns them alongside the repos (needed by main for audit middleware).
func New(db *pgxpool.Pool, cfg *config.Config) (*Handlers, *repository.Repos) {
	repos := repository.New(db)
	initial, err := repos.Operator.LoadState(context.Background())
	if err != nil {
		initial = maintenance.Snapshot{}
	}
	state := maintenance.New(initial)

	authSvc := services.NewAuthService(repos.Users, repos.Tokens, cfg)
	userSvc := services.NewUserService(repos.Users, repos.Roles, repos.Tokens)
	appSvc := services.NewApplicationService(repos.Applications, repos.Audit, repos.Organisations)
	ussdSvc := services.NewUSSDService(repos.USSD, repos.Applications)
	formSvc := services.NewFormService(repos.Forms, repos.Departments, repos.Audit)
	secSvc := services.NewSecurityService(repos.Security, repos.Users)
	updSvc := services.NewUpdatesService(db)
	deployer := services.NewDeployer()

	return &Handlers{
		Auth:          &AuthHandler{svc: authSvc, users: repos.Users, cfg: cfg},
		Users:         &UsersHandler{svc: userSvc, audit: repos.Audit},
		Roles:         &RolesHandler{roles: repos.Roles, audit: repos.Audit},
		Applications:  &ApplicationsHandler{svc: appSvc, audit: repos.Audit, storage: repos.CMS},
		Dashboard:     &DashboardHandler{users: repos.Users, apps: repos.Applications},
		Audit:         &AuditHandler{repo: repos.Audit},
		CMS:           &CMSHandler{repo: repos.CMS},
		NSMIS:         &NSMISHandler{repo: repos.NSMIS, cfg: cfg},
		USSD:          &USSDHandler{service: ussdSvc, callbackSecret: cfg.USSDCallbackSecret},
		Organisations: &OrganisationsHandler{repo: repos.Organisations},
		Operator:      NewOperatorHandler(repos.Operator, state),
		SystemState:   state,
		Backups:       &BackupsHandler{repo: repos.Backups, cfg: cfg},
		Forms:         &FormsHandler{svc: formSvc},
		Security:      &SecurityHandler{svc: secSvc, users: repos.Users, audit: repos.Audit},
		Updates:       &UpdatesHandler{svc: updSvc, deployer: deployer, backups: repos.Backups},
		UpdatesSvc:    updSvc,
		Analytics:     NewAnalyticsHandler(repos.Analytics, cfg),
	}, repos
}
