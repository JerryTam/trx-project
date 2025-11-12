package service

import (
	"context"
	"errors"
	"trx-project/internal/model"
	"trx-project/internal/repository"
	"trx-project/pkg/cache"

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
	repo        repository.RBACRepository
	cache       *cache.RBACCache
	logger      *zap.Logger
	enableCache bool // 是否启用缓存
}

// NewRBACService 创建 RBAC 服务
func NewRBACService(repo repository.RBACRepository, rbacCache *cache.RBACCache, logger *zap.Logger) RBACService {
	return &rbacService{
		repo:        repo,
		cache:       rbacCache,
		logger:      logger,
		enableCache: rbacCache != nil, // 如果提供了缓存，则启用
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
	err := s.repo.AssignPermissionsToRole(ctx, roleID, permissionIDs)
	if err != nil {
		return err
	}

	// 使角色权限缓存失效
	if s.enableCache {
		if err := s.cache.InvalidateRoleCache(ctx, roleID); err != nil {
			s.logger.Error("Failed to invalidate role cache",
				zap.Uint("role_id", roleID),
				zap.Error(err))
		}
	}

	return nil
}

func (s *rbacService) RemovePermissionsFromRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	err := s.repo.RemovePermissionsFromRole(ctx, roleID, permissionIDs)
	if err != nil {
		return err
	}

	// 使角色权限缓存失效
	if s.enableCache {
		if err := s.cache.InvalidateRoleCache(ctx, roleID); err != nil {
			s.logger.Error("Failed to invalidate role cache",
				zap.Uint("role_id", roleID),
				zap.Error(err))
		}
	}

	return nil
}

func (s *rbacService) GetRolePermissions(ctx context.Context, roleID uint) ([]*model.Permission, error) {
	// 尝试从缓存获取
	if s.enableCache {
		if permissionCodes, ok := s.cache.GetRolePermissions(ctx, roleID); ok {
			// 缓存命中，转换为 Permission 对象
			permissions := make([]*model.Permission, 0, len(permissionCodes))
			for _, code := range permissionCodes {
				permissions = append(permissions, &model.Permission{Code: code})
			}
			return permissions, nil
		}
	}

	// 缓存未命中，从数据库获取
	permissions, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 存入缓存
	if s.enableCache && len(permissions) > 0 {
		permissionCodes := make([]string, 0, len(permissions))
		for _, p := range permissions {
			permissionCodes = append(permissionCodes, p.Code)
		}
		if err := s.cache.SetRolePermissions(ctx, roleID, permissionCodes); err != nil {
			s.logger.Error("Failed to cache role permissions",
				zap.Uint("role_id", roleID),
				zap.Error(err))
		}
	}

	return permissions, nil
}

// UserRole 相关实现

func (s *rbacService) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	// 检查角色是否存在
	_, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return errors.New("role not found")
	}

	err = s.repo.AssignRoleToUser(ctx, userID, roleID)
	if err != nil {
		return err
	}

	// 使用户缓存失效
	if s.enableCache {
		if err := s.cache.InvalidateUserCache(ctx, userID); err != nil {
			s.logger.Error("Failed to invalidate user cache",
				zap.Uint("user_id", userID),
				zap.Error(err))
		}
	}

	return nil
}

func (s *rbacService) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	err := s.repo.RemoveRoleFromUser(ctx, userID, roleID)
	if err != nil {
		return err
	}

	// 使用户缓存失效
	if s.enableCache {
		if err := s.cache.InvalidateUserCache(ctx, userID); err != nil {
			s.logger.Error("Failed to invalidate user cache",
				zap.Uint("user_id", userID),
				zap.Error(err))
		}
	}

	return nil
}

func (s *rbacService) GetUserRoles(ctx context.Context, userID uint) ([]*model.Role, error) {
	return s.repo.GetUserRoles(ctx, userID)
}

func (s *rbacService) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	// 尝试从缓存获取
	if s.enableCache {
		if permissionCodes, ok := s.cache.GetUserPermissions(ctx, userID); ok {
			// 缓存命中，转换为 Permission 对象
			permissions := make([]*model.Permission, 0, len(permissionCodes))
			for _, code := range permissionCodes {
				permissions = append(permissions, &model.Permission{Code: code})
			}
			s.logger.Debug("User permissions cache hit", zap.Uint("user_id", userID))
			return permissions, nil
		}
	}

	// 缓存未命中，从数据库获取
	permissions, err := s.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 存入缓存
	if s.enableCache && len(permissions) > 0 {
		permissionCodes := make([]string, 0, len(permissions))
		for _, p := range permissions {
			permissionCodes = append(permissionCodes, p.Code)
		}
		if err := s.cache.SetUserPermissions(ctx, userID, permissionCodes); err != nil {
			s.logger.Error("Failed to cache user permissions",
				zap.Uint("user_id", userID),
				zap.Error(err))
		}
	}

	return permissions, nil
}

func (s *rbacService) HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error) {
	// 尝试从缓存获取权限检查结果
	if s.enableCache {
		if hasPermission, ok := s.cache.CheckPermissionCached(ctx, userID, permissionCode); ok {
			// 缓存命中
			s.logger.Debug("Permission check cache hit",
				zap.Uint("user_id", userID),
				zap.String("permission", permissionCode),
				zap.Bool("has_permission", hasPermission))
			return hasPermission, nil
		}
	}

	// 缓存未命中，从数据库查询
	hasPermission, err := s.repo.HasPermission(ctx, userID, permissionCode)
	if err != nil {
		return false, err
	}

	// 存入缓存
	if s.enableCache {
		if err := s.cache.SetPermissionCheck(ctx, userID, permissionCode, hasPermission); err != nil {
			s.logger.Error("Failed to cache permission check",
				zap.Uint("user_id", userID),
				zap.String("permission", permissionCode),
				zap.Error(err))
		}
	}

	return hasPermission, nil
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
