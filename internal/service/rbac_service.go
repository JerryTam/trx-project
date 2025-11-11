package service

import (
	"context"
	"errors"
	"trx-project/internal/model"
	"trx-project/internal/repository"

	"go.uber.org/zap"
)

// RBACService RBAC 服务接口
type RBACService interface {
	// Role 相关
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetRoleWithPermissions(ctx context.Context, roleID uint) (*model.Role, error)
	ListRoles(ctx context.Context) ([]*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) error
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, id uint) error

	// Permission 相关
	ListPermissions(ctx context.Context) ([]*model.Permission, error)
	CreatePermission(ctx context.Context, permission *model.Permission) error

	// RolePermission 相关
	AssignPermissionsToRole(ctx context.Context, roleID uint, permissionIDs []uint) error
	RemovePermissionsFromRole(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]*model.Permission, error)

	// UserRole 相关
	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]*model.Role, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error)
	HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error)
	CheckPermission(ctx context.Context, userID uint, permissionCode string) error
}

type rbacService struct {
	repo   repository.RBACRepository
	logger *zap.Logger
}

// NewRBACService 创建 RBAC 服务
func NewRBACService(repo repository.RBACRepository, logger *zap.Logger) RBACService {
	return &rbacService{
		repo:   repo,
		logger: logger,
	}
}

// Role 相关实现

func (s *rbacService) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	return s.repo.GetRoleByName(ctx, name)
}

func (s *rbacService) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	return s.repo.GetRoleByID(ctx, id)
}

func (s *rbacService) GetRoleWithPermissions(ctx context.Context, roleID uint) (*model.Role, error) {
	return s.repo.GetRoleWithPermissions(ctx, roleID)
}

func (s *rbacService) ListRoles(ctx context.Context) ([]*model.Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *rbacService) CreateRole(ctx context.Context, role *model.Role) error {
	// 检查角色名是否已存在
	existingRole, err := s.repo.GetRoleByName(ctx, role.Name)
	if err == nil && existingRole.ID > 0 {
		return errors.New("role name already exists")
	}
	
	return s.repo.CreateRole(ctx, role)
}

func (s *rbacService) UpdateRole(ctx context.Context, role *model.Role) error {
	return s.repo.UpdateRole(ctx, role)
}

func (s *rbacService) DeleteRole(ctx context.Context, id uint) error {
	return s.repo.DeleteRole(ctx, id)
}

// Permission 相关实现

func (s *rbacService) ListPermissions(ctx context.Context) ([]*model.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *rbacService) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return s.repo.CreatePermission(ctx, permission)
}

// RolePermission 相关实现

func (s *rbacService) AssignPermissionsToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return s.repo.AssignPermissionsToRole(ctx, roleID, permissionIDs)
}

func (s *rbacService) RemovePermissionsFromRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return s.repo.RemovePermissionsFromRole(ctx, roleID, permissionIDs)
}

func (s *rbacService) GetRolePermissions(ctx context.Context, roleID uint) ([]*model.Permission, error) {
	return s.repo.GetRolePermissions(ctx, roleID)
}

// UserRole 相关实现

func (s *rbacService) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	// 检查角色是否存在
	_, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return errors.New("role not found")
	}
	
	return s.repo.AssignRoleToUser(ctx, userID, roleID)
}

func (s *rbacService) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	return s.repo.RemoveRoleFromUser(ctx, userID, roleID)
}

func (s *rbacService) GetUserRoles(ctx context.Context, userID uint) ([]*model.Role, error) {
	return s.repo.GetUserRoles(ctx, userID)
}

func (s *rbacService) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	return s.repo.GetUserPermissions(ctx, userID)
}

func (s *rbacService) HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error) {
	return s.repo.HasPermission(ctx, userID, permissionCode)
}

// CheckPermission 检查用户是否有指定权限，没有则返回错误
func (s *rbacService) CheckPermission(ctx context.Context, userID uint, permissionCode string) error {
	has, err := s.repo.HasPermission(ctx, userID, permissionCode)
	if err != nil {
		s.logger.Error("Failed to check permission", 
			zap.Uint("user_id", userID), 
			zap.String("permission", permissionCode),
			zap.Error(err))
		return err
	}
	
	if !has {
		s.logger.Warn("Permission denied", 
			zap.Uint("user_id", userID), 
			zap.String("permission", permissionCode))
		return errors.New("permission denied")
	}
	
	return nil
}

