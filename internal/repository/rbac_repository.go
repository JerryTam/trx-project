package repository

import (
	"context"
	"trx-project/internal/model"

	"gorm.io/gorm"
)

// RBACRepository RBAC 数据访问接口
type RBACRepository interface {
	// Role 相关
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetRoleWithPermissions(ctx context.Context, roleID uint) (*model.Role, error)
	ListRoles(ctx context.Context) ([]*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role) error
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, id uint) error

	// Permission 相关
	GetPermissionByCode(ctx context.Context, code string) (*model.Permission, error)
	GetPermissionByID(ctx context.Context, id uint) (*model.Permission, error)
	ListPermissions(ctx context.Context) ([]*model.Permission, error)
	CreatePermission(ctx context.Context, permission *model.Permission) error
	UpdatePermission(ctx context.Context, permission *model.Permission) error
	DeletePermission(ctx context.Context, id uint) error

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
}

type rbacRepository struct {
	db *gorm.DB
}

// NewRBACRepository 创建 RBAC repository
func NewRBACRepository(db *gorm.DB) RBACRepository {
	return &rbacRepository{db: db}
}

// Role 相关实现

func (r *rbacRepository) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *rbacRepository) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	return &role, err
}

func (r *rbacRepository) GetRoleWithPermissions(ctx context.Context, roleID uint) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Preload("Permissions").First(&role, roleID).Error
	return &role, err
}

func (r *rbacRepository) ListRoles(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}

func (r *rbacRepository) CreateRole(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *rbacRepository) UpdateRole(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *rbacRepository) DeleteRole(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// Permission 相关实现

func (r *rbacRepository) GetPermissionByCode(ctx context.Context, code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&permission).Error
	return &permission, err
}

func (r *rbacRepository) GetPermissionByID(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).First(&permission, id).Error
	return &permission, err
}

func (r *rbacRepository) ListPermissions(ctx context.Context) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	return permissions, err
}

func (r *rbacRepository) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *rbacRepository) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

func (r *rbacRepository) DeletePermission(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Permission{}, id).Error
}

// RolePermission 相关实现

func (r *rbacRepository) AssignPermissionsToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	role := model.Role{ID: roleID}

	var permissions []model.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Append(permissions)
}

func (r *rbacRepository) RemovePermissionsFromRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	role := model.Role{ID: roleID}

	var permissions []model.Permission
	if err := r.db.WithContext(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Delete(permissions)
}

func (r *rbacRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// UserRole 相关实现

func (r *rbacRepository) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	userRole := &model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *rbacRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{}).Error
}

func (r *rbacRepository) GetUserRoles(ctx context.Context, userID uint) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *rbacRepository) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).
		Distinct().
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.status = 1", userID).
		Find(&permissions).Error
	return permissions, err
}

func (r *rbacRepository) HasPermission(ctx context.Context, userID uint, permissionCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.code = ? AND permissions.status = 1", userID, permissionCode).
		Count(&count).Error

	return count > 0, err
}
