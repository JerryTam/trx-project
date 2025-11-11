package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:50" json:"name"`        // 角色名称：superadmin, admin, editor, viewer
	DisplayName string         `gorm:"not null;size:100" json:"display_name"`           // 显示名称：超级管理员、管理员、编辑、查看者
	Description string         `gorm:"size:500" json:"description"`                     // 角色描述
	Status      int            `gorm:"default:1;not null" json:"status"`                // 状态：1-启用 0-禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"` // 角色拥有的权限
}

// Permission 权限模型
type Permission struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Code        string         `gorm:"uniqueIndex;not null;size:100" json:"code"`  // 权限编码：user:read, user:write
	Name        string         `gorm:"not null;size:100" json:"name"`              // 权限名称：查看用户、编辑用户
	Resource    string         `gorm:"not null;size:50" json:"resource"`           // 资源：user, role, permission, statistics
	Action      string         `gorm:"not null;size:50" json:"action"`             // 操作：read, write, delete
	Description string         `gorm:"size:500" json:"description"`                // 权限描述
	Status      int            `gorm:"default:1;not null" json:"status"`           // 状态：1-启用 0-禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint      `gorm:"primarykey" json:"role_id"`
	PermissionID uint      `gorm:"primarykey" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserRole 用户角色关联模型
type UserRole struct {
	UserID    uint      `gorm:"primarykey" json:"user_id"`
	RoleID    uint      `gorm:"primarykey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (UserRole) TableName() string {
	return "user_roles"
}

