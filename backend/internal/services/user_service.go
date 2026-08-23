package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/google/uuid"
)

// ErrProtectedAccount is returned when an operation targets the protected super admin.
var ErrProtectedAccount = errors.New("this account is protected and cannot be modified or deleted")
var ErrSelfAccountAction = errors.New("you cannot apply this account action to your own account")
var ErrFraudFlagged = errors.New("clear the fraud flag before reactivating this account")
var ErrInvalidAccountAction = errors.New("invalid account action")

type UserService struct {
	users  *repository.UserRepo
	roles  *repository.RoleRepo
	tokens *repository.TokenRepo
}

func NewUserService(users *repository.UserRepo, roles *repository.RoleRepo, tokens *repository.TokenRepo) *UserService {
	return &UserService{users: users, roles: roles, tokens: tokens}
}

type CreateUserInput struct {
	Email        string
	Password     string
	FirstName    string
	LastName     string
	Phone        string
	NIN          string
	FederationID string
	ActorID      string
}

func (s *UserService) Create(ctx context.Context, in CreateUserInput) (*models.User, error) {
	hash, err := HashPassword(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &models.User{
		ID:            uuid.NewString(),
		Email:         in.Email,
		PasswordHash:  hash,
		FirstName:     in.FirstName,
		LastName:      in.LastName,
		Phone:         in.Phone,
		NIN:           in.NIN,
		IsActive:      true,
		AccountStatus: models.AccountStatusActive,
	}

	if err := s.users.Create(ctx, u); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, fmt.Errorf("email already registered")
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	if in.FederationID != "" {
		if err := s.users.CreateFederationMembership(ctx, u.ID, in.FederationID, "OFFICER", in.ActorID); err != nil {
			return nil, fmt.Errorf("link federation: %w", err)
		}
	}

	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	roles, err := s.users.GetRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Roles = roles
	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) List(ctx context.Context, p *models.PaginationParams, federationIDs []string) ([]*models.User, int64, error) {
	users, total, err := s.users.List(ctx, p, federationIDs)
	if err != nil {
		return nil, 0, err
	}
	for _, u := range users {
		u.PasswordHash = ""
		roles, _ := s.users.GetRoles(ctx, u.ID)
		u.Roles = roles
	}
	return users, total, nil
}

func (s *UserService) HasPermission(ctx context.Context, userID string, perms ...string) (bool, error) {
	return s.users.HasAnyPermission(ctx, userID, perms...)
}

func (s *UserService) GetFederationIDs(ctx context.Context, userID string) ([]string, error) {
	return s.users.GetFederationIDs(ctx, userID)
}

func (s *UserService) IsUserInFederations(ctx context.Context, targetUserID string, federationIDs []string) (bool, error) {
	return s.users.IsUserInFederations(ctx, targetUserID, federationIDs)
}

func (s *UserService) LinkFederation(ctx context.Context, userID, federationID, actorID string) error {
	return s.users.CreateFederationMembership(ctx, userID, federationID, "OFFICER", actorID)
}

type UpdateUserInput struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	NIN       string
}

func (s *UserService) Update(ctx context.Context, id string, in UpdateUserInput) (*models.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.FirstName = in.FirstName
	user.LastName = in.LastName
	if in.Email != "" {
		user.Email = in.Email
	}
	user.Phone = in.Phone
	if in.NIN != "" {
		user.NIN = in.NIN
	}
	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if err := s.guardSuperAdmin(ctx, id); err != nil {
		return err
	}
	if err := s.users.SoftDelete(ctx, id); err != nil {
		return err
	}
	return s.tokens.RevokeAllForUser(ctx, id)
}

func (s *UserService) SetActive(ctx context.Context, id string, active bool) error {
	if !active {
		if err := s.guardSuperAdmin(ctx, id); err != nil {
			return err
		}
	}
	if err := s.users.SetActive(ctx, id, active); err != nil {
		return err
	}
	if !active {
		return s.tokens.RevokeAllForUser(ctx, id)
	}
	return nil
}

var sportsRegistryRoles = map[string]bool{
	"club_manager":        true,
	"coach":               true,
	"athlete":             true,
	"technical_official":  true,
	"medical_officer":     true,
	"anti_doping_officer": true,
	"federation_officer":  true,
}

func isSportsRegistryRole(roleName string) bool {
	return sportsRegistryRoles[roleName]
}

func (s *UserService) AssignRole(ctx context.Context, userID, roleID, assignedBy string) error {
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return repository.ErrNotFound
	}
	role, err := s.roles.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found")
	}

	// Permission scoping check
	globalAssign, err := s.users.HasAnyPermission(ctx, assignedBy, "users:roles", "users:write", "*", "users:assign")
	if err != nil {
		return fmt.Errorf("permission check: %w", err)
	}
	if !globalAssign {
		roles, rErr := s.users.GetRoles(ctx, assignedBy)
		if rErr == nil {
			for _, r := range roles {
				if r.Name == "super_admin" || r.Name == "admin" {
					globalAssign = true
					break
				}
			}
		}
	}

	if !globalAssign {
		scopedAssign, err := s.users.HasAnyPermission(ctx, assignedBy, "users:roles:own")
		if err != nil || !scopedAssign {
			return errors.New("unauthorized to assign roles")
		}

		// Verify target user is in assigner's federation
		federations, err := s.users.GetFederationIDs(ctx, assignedBy)
		if err != nil {
			return fmt.Errorf("get federations: %w", err)
		}
		inFed, err := s.users.IsUserInFederations(ctx, userID, federations)
		if err != nil || !inFed {
			return errors.New("user does not belong to your federation")
		}

		// Verify role is a sports registry role
		if !isSportsRegistryRole(role.Name) {
			return fmt.Errorf("role %s cannot be assigned by federation profile", role.Name)
		}
	}

	// Enforce single super_admin — block a second assignment of the role
	if role.Name == "super_admin" {
		count, err := s.users.CountSuperAdmins(ctx)
		if err != nil {
			return fmt.Errorf("check super admin count: %w", err)
		}
		if count >= 1 {
			return ErrProtectedAccount
		}
	}
	return s.users.AssignRole(ctx, userID, roleID, assignedBy)
}

func (s *UserService) RemoveRole(ctx context.Context, userID, roleID, actorID string) error {
	role, err := s.roles.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found")
	}

	// Permission scoping check
	globalAssign, err := s.users.HasAnyPermission(ctx, actorID, "users:roles", "users:write", "*", "users:assign")
	if err != nil {
		return fmt.Errorf("permission check: %w", err)
	}
	if !globalAssign {
		roles, rErr := s.users.GetRoles(ctx, actorID)
		if rErr == nil {
			for _, r := range roles {
				if r.Name == "super_admin" || r.Name == "admin" {
					globalAssign = true
					break
				}
			}
		}
	}

	if !globalAssign {
		scopedAssign, err := s.users.HasAnyPermission(ctx, actorID, "users:roles:own")
		if err != nil || !scopedAssign {
			return errors.New("unauthorized to remove roles")
		}

		// Verify target user is in assigner's federation
		federations, err := s.users.GetFederationIDs(ctx, actorID)
		if err != nil {
			return fmt.Errorf("get federations: %w", err)
		}
		inFed, err := s.users.IsUserInFederations(ctx, userID, federations)
		if err != nil || !inFed {
			return errors.New("user does not belong to your federation")
		}

		// Verify role is a sports registry role
		if !isSportsRegistryRole(role.Name) {
			return fmt.Errorf("role %s cannot be managed by federation profile", role.Name)
		}
	}

	if role.Name == "super_admin" {
		return ErrProtectedAccount
	}
	return s.users.RemoveRole(ctx, userID, roleID)
}

func (s *UserService) ResetPassword(ctx context.Context, userID, newPassword string) error {
	if err := s.guardSuperAdmin(ctx, userID); err != nil {
		return err
	}
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return err
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err := s.users.UpdatePassword(ctx, userID, hash); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if err := s.tokens.RevokeAllForUser(ctx, userID); err != nil {
		return fmt.Errorf("password updated but session revocation failed: %w", err)
	}
	return nil
}

type AccountActionInput struct {
	UserID         string
	ActorID        string
	Action         string
	Reason         string
	SuspendedUntil *time.Time
}

func (s *UserService) ApplyAccountAction(ctx context.Context, in AccountActionInput) error {
	if in.UserID == in.ActorID {
		return ErrSelfAccountAction
	}
	if err := s.guardSuperAdmin(ctx, in.UserID); err != nil {
		return err
	}
	user, err := s.users.GetByID(ctx, in.UserID)
	if err != nil {
		return err
	}

	action := strings.ToUpper(strings.TrimSpace(in.Action))
	reason := strings.TrimSpace(in.Reason)
	status := user.AccountStatus
	fraud := user.FraudFlag
	fraudReason := user.FraudReason
	suspendedUntil := user.SuspendedUntil
	restrictAccess := false

	switch action {
	case "SUSPEND":
		if reason == "" {
			return errors.New("a reason is required to suspend an account")
		}
		status, suspendedUntil, restrictAccess = models.AccountStatusSuspended, in.SuspendedUntil, true
	case "BAN":
		if reason == "" {
			return errors.New("a reason is required to ban an account")
		}
		status, suspendedUntil, restrictAccess = models.AccountStatusBanned, nil, true
	case "MARK_FRAUD":
		if reason == "" {
			return errors.New("a reason is required to mark an account as fraud")
		}
		status, fraud, fraudReason, suspendedUntil, restrictAccess = models.AccountStatusBanned, true, reason, nil, true
	case "CLEAR_FRAUD":
		fraud, fraudReason = false, ""
	case "REACTIVATE":
		if fraud {
			return ErrFraudFlagged
		}
		status, reason, suspendedUntil = models.AccountStatusActive, "", nil
	default:
		return ErrInvalidAccountAction
	}

	if err := s.users.SetAccountManagement(ctx, in.UserID, in.ActorID, status, reason, fraud, fraudReason, suspendedUntil); err != nil {
		return fmt.Errorf("apply account action: %w", err)
	}
	if restrictAccess {
		if err := s.tokens.RevokeAllForUser(ctx, in.UserID); err != nil {
			return fmt.Errorf("account restricted but session revocation failed: %w", err)
		}
	}
	return nil
}

// guardSuperAdmin rejects any mutating operation targeting a user with the super_admin role.
func (s *UserService) guardSuperAdmin(ctx context.Context, userID string) error {
	roles, err := s.users.GetRoles(ctx, userID)
	if err != nil {
		return err
	}
	for _, r := range roles {
		if r.Name == "super_admin" {
			return ErrProtectedAccount
		}
	}
	return nil
}
