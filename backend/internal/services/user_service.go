package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/google/uuid"
)

// ErrProtectedAccount is returned when an operation targets the protected super admin.
var ErrProtectedAccount = errors.New("this account is protected and cannot be modified or deleted")

type UserService struct {
	users *repository.UserRepo
	roles *repository.RoleRepo
}

func NewUserService(users *repository.UserRepo, roles *repository.RoleRepo) *UserService {
	return &UserService{users: users, roles: roles}
}

type CreateUserInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
}

func (s *UserService) Create(ctx context.Context, in CreateUserInput) (*models.User, error) {
	hash, err := HashPassword(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &models.User{
		ID:           uuid.NewString(),
		Email:        in.Email,
		PasswordHash: hash,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Phone:        in.Phone,
		IsActive:     true,
	}

	if err := s.users.Create(ctx, u); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, fmt.Errorf("email already registered")
		}
		return nil, fmt.Errorf("create user: %w", err)
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

func (s *UserService) List(ctx context.Context, p *models.PaginationParams) ([]*models.User, int64, error) {
	users, total, err := s.users.List(ctx, p)
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

type UpdateUserInput struct {
	FirstName string
	LastName  string
	Phone     string
}

func (s *UserService) Update(ctx context.Context, id string, in UpdateUserInput) (*models.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.FirstName = in.FirstName
	user.LastName = in.LastName
	user.Phone = in.Phone
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
	return s.users.SoftDelete(ctx, id)
}

func (s *UserService) SetActive(ctx context.Context, id string, active bool) error {
	if !active {
		if err := s.guardSuperAdmin(ctx, id); err != nil {
			return err
		}
	}
	return s.users.SetActive(ctx, id, active)
}

func (s *UserService) AssignRole(ctx context.Context, userID, roleID, assignedBy string) error {
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return repository.ErrNotFound
	}
	role, err := s.roles.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found")
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

func (s *UserService) RemoveRole(ctx context.Context, userID, roleID string) error {
	role, err := s.roles.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found")
	}
	if role.Name == "super_admin" {
		return ErrProtectedAccount
	}
	return s.users.RemoveRole(ctx, userID, roleID)
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
