package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RBACCache RBAC 权限缓存管理
type RBACCache struct {
	redis  *redis.Client
	logger *zap.Logger
	// 缓存过期时间
	userRolesTTL        time.Duration // 用户角色缓存
	rolePermissionsTTL  time.Duration // 角色权限缓存
	userPermissionsTTL  time.Duration // 用户权限缓存（聚合）
}

// NewRBACCache 创建 RBAC 缓存管理器
func NewRBACCache(redis *redis.Client, logger *zap.Logger) *RBACCache {
	return &RBACCache{
		redis:               redis,
		logger:              logger,
		userRolesTTL:        5 * time.Minute,  // 用户角色缓存 5 分钟
		rolePermissionsTTL:  10 * time.Minute, // 角色权限缓存 10 分钟
		userPermissionsTTL:  5 * time.Minute,  // 用户权限缓存 5 分钟
	}
}

// Cache Keys 定义
const (
	// 用户角色列表: rbac:user_roles:<user_id>
	userRolesKeyPrefix = "rbac:user_roles:"
	// 角色权限列表: rbac:role_permissions:<role_id>
	rolePermissionsKeyPrefix = "rbac:role_permissions:"
	// 用户权限列表（聚合）: rbac:user_permissions:<user_id>
	userPermissionsKeyPrefix = "rbac:user_permissions:"
	// 权限检查结果: rbac:check:<user_id>:<permission>
	permissionCheckKeyPrefix = "rbac:check:"
)

// === 用户角色缓存 ===

// GetUserRoles 获取用户角色（带缓存）
func (c *RBACCache) GetUserRoles(ctx context.Context, userID uint) ([]uint, bool) {
	key := fmt.Sprintf("%s%d", userRolesKeyPrefix, userID)
	
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			c.logger.Error("Failed to get user roles from cache",
				zap.Uint("user_id", userID),
				zap.Error(err))
		}
		return nil, false
	}
	
	var roleIDs []uint
	if err := json.Unmarshal([]byte(data), &roleIDs); err != nil {
		c.logger.Error("Failed to unmarshal user roles",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return nil, false
	}
	
	c.logger.Debug("User roles cache hit",
		zap.Uint("user_id", userID),
		zap.Int("role_count", len(roleIDs)))
	
	return roleIDs, true
}

// SetUserRoles 设置用户角色缓存
func (c *RBACCache) SetUserRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	key := fmt.Sprintf("%s%d", userRolesKeyPrefix, userID)
	
	data, err := json.Marshal(roleIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal user roles: %w", err)
	}
	
	if err := c.redis.Set(ctx, key, data, c.userRolesTTL).Err(); err != nil {
		return fmt.Errorf("failed to set user roles cache: %w", err)
	}
	
	c.logger.Debug("User roles cached",
		zap.Uint("user_id", userID),
		zap.Int("role_count", len(roleIDs)),
		zap.Duration("ttl", c.userRolesTTL))
	
	return nil
}

// DeleteUserRoles 删除用户角色缓存
func (c *RBACCache) DeleteUserRoles(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("%s%d", userRolesKeyPrefix, userID)
	
	if err := c.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete user roles cache: %w", err)
	}
	
	c.logger.Debug("User roles cache deleted", zap.Uint("user_id", userID))
	return nil
}

// === 角色权限缓存 ===

// GetRolePermissions 获取角色权限（带缓存）
func (c *RBACCache) GetRolePermissions(ctx context.Context, roleID uint) ([]string, bool) {
	key := fmt.Sprintf("%s%d", rolePermissionsKeyPrefix, roleID)
	
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			c.logger.Error("Failed to get role permissions from cache",
				zap.Uint("role_id", roleID),
				zap.Error(err))
		}
		return nil, false
	}
	
	var permissions []string
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		c.logger.Error("Failed to unmarshal role permissions",
			zap.Uint("role_id", roleID),
			zap.Error(err))
		return nil, false
	}
	
	c.logger.Debug("Role permissions cache hit",
		zap.Uint("role_id", roleID),
		zap.Int("permission_count", len(permissions)))
	
	return permissions, true
}

// SetRolePermissions 设置角色权限缓存
func (c *RBACCache) SetRolePermissions(ctx context.Context, roleID uint, permissions []string) error {
	key := fmt.Sprintf("%s%d", rolePermissionsKeyPrefix, roleID)
	
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal role permissions: %w", err)
	}
	
	if err := c.redis.Set(ctx, key, data, c.rolePermissionsTTL).Err(); err != nil {
		return fmt.Errorf("failed to set role permissions cache: %w", err)
	}
	
	c.logger.Debug("Role permissions cached",
		zap.Uint("role_id", roleID),
		zap.Int("permission_count", len(permissions)),
		zap.Duration("ttl", c.rolePermissionsTTL))
	
	return nil
}

// DeleteRolePermissions 删除角色权限缓存
func (c *RBACCache) DeleteRolePermissions(ctx context.Context, roleID uint) error {
	key := fmt.Sprintf("%s%d", rolePermissionsKeyPrefix, roleID)
	
	if err := c.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete role permissions cache: %w", err)
	}
	
	c.logger.Debug("Role permissions cache deleted", zap.Uint("role_id", roleID))
	return nil
}

// === 用户权限缓存（聚合） ===

// GetUserPermissions 获取用户所有权限（带缓存）
func (c *RBACCache) GetUserPermissions(ctx context.Context, userID uint) ([]string, bool) {
	key := fmt.Sprintf("%s%d", userPermissionsKeyPrefix, userID)
	
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			c.logger.Error("Failed to get user permissions from cache",
				zap.Uint("user_id", userID),
				zap.Error(err))
		}
		return nil, false
	}
	
	var permissions []string
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		c.logger.Error("Failed to unmarshal user permissions",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return nil, false
	}
	
	c.logger.Debug("User permissions cache hit",
		zap.Uint("user_id", userID),
		zap.Int("permission_count", len(permissions)))
	
	return permissions, true
}

// SetUserPermissions 设置用户权限缓存
func (c *RBACCache) SetUserPermissions(ctx context.Context, userID uint, permissions []string) error {
	key := fmt.Sprintf("%s%d", userPermissionsKeyPrefix, userID)
	
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal user permissions: %w", err)
	}
	
	if err := c.redis.Set(ctx, key, data, c.userPermissionsTTL).Err(); err != nil {
		return fmt.Errorf("failed to set user permissions cache: %w", err)
	}
	
	c.logger.Debug("User permissions cached",
		zap.Uint("user_id", userID),
		zap.Int("permission_count", len(permissions)),
		zap.Duration("ttl", c.userPermissionsTTL))
	
	return nil
}

// DeleteUserPermissions 删除用户权限缓存
func (c *RBACCache) DeleteUserPermissions(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("%s%d", userPermissionsKeyPrefix, userID)
	
	if err := c.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete user permissions cache: %w", err)
	}
	
	c.logger.Debug("User permissions cache deleted", zap.Uint("user_id", userID))
	return nil
}

// === 权限检查缓存 ===

// CheckPermissionCached 检查权限（带缓存）
func (c *RBACCache) CheckPermissionCached(ctx context.Context, userID uint, permission string) (bool, bool) {
	key := fmt.Sprintf("%s%d:%s", permissionCheckKeyPrefix, userID, permission)
	
	result, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			c.logger.Error("Failed to get permission check from cache",
				zap.Uint("user_id", userID),
				zap.String("permission", permission),
				zap.Error(err))
		}
		return false, false
	}
	
	hasPermission := result == "1"
	
	c.logger.Debug("Permission check cache hit",
		zap.Uint("user_id", userID),
		zap.String("permission", permission),
		zap.Bool("has_permission", hasPermission))
	
	return hasPermission, true
}

// SetPermissionCheck 设置权限检查结果缓存
func (c *RBACCache) SetPermissionCheck(ctx context.Context, userID uint, permission string, hasPermission bool) error {
	key := fmt.Sprintf("%s%d:%s", permissionCheckKeyPrefix, userID, permission)
	
	value := "0"
	if hasPermission {
		value = "1"
	}
	
	if err := c.redis.Set(ctx, key, value, c.userPermissionsTTL).Err(); err != nil {
		return fmt.Errorf("failed to set permission check cache: %w", err)
	}
	
	c.logger.Debug("Permission check cached",
		zap.Uint("user_id", userID),
		zap.String("permission", permission),
		zap.Bool("has_permission", hasPermission),
		zap.Duration("ttl", c.userPermissionsTTL))
	
	return nil
}

// DeletePermissionChecks 删除用户的所有权限检查缓存
func (c *RBACCache) DeletePermissionChecks(ctx context.Context, userID uint) error {
	pattern := fmt.Sprintf("%s%d:*", permissionCheckKeyPrefix, userID)
	
	// 使用 SCAN 命令查找所有匹配的 key
	var cursor uint64
	var keys []string
	
	for {
		var batch []string
		var err error
		batch, cursor, err = c.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan permission check keys: %w", err)
		}
		
		keys = append(keys, batch...)
		
		if cursor == 0 {
			break
		}
	}
	
	// 批量删除
	if len(keys) > 0 {
		if err := c.redis.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete permission check caches: %w", err)
		}
		
		c.logger.Debug("Permission check caches deleted",
			zap.Uint("user_id", userID),
			zap.Int("count", len(keys)))
	}
	
	return nil
}

// === 批量缓存失效 ===

// InvalidateUserCache 使用户的所有缓存失效
func (c *RBACCache) InvalidateUserCache(ctx context.Context, userID uint) error {
	// 删除用户角色缓存
	if err := c.DeleteUserRoles(ctx, userID); err != nil {
		c.logger.Error("Failed to delete user roles cache", zap.Error(err))
	}
	
	// 删除用户权限缓存
	if err := c.DeleteUserPermissions(ctx, userID); err != nil {
		c.logger.Error("Failed to delete user permissions cache", zap.Error(err))
	}
	
	// 删除权限检查缓存
	if err := c.DeletePermissionChecks(ctx, userID); err != nil {
		c.logger.Error("Failed to delete permission checks cache", zap.Error(err))
	}
	
	c.logger.Info("User cache invalidated", zap.Uint("user_id", userID))
	return nil
}

// InvalidateRoleCache 使角色的所有缓存失效
func (c *RBACCache) InvalidateRoleCache(ctx context.Context, roleID uint) error {
	// 删除角色权限缓存
	if err := c.DeleteRolePermissions(ctx, roleID); err != nil {
		c.logger.Error("Failed to delete role permissions cache", zap.Error(err))
		return err
	}
	
	c.logger.Info("Role cache invalidated", zap.Uint("role_id", roleID))
	
	// 注意：这里不删除用户缓存，因为我们不知道哪些用户拥有这个角色
	// 用户缓存会在 TTL 过期后自动更新
	// 如果需要立即生效，可以在分配/移除角色时手动调用 InvalidateUserCache
	
	return nil
}

// InvalidateAllRBACCache 清除所有 RBAC 缓存（谨慎使用）
func (c *RBACCache) InvalidateAllRBACCache(ctx context.Context) error {
	patterns := []string{
		userRolesKeyPrefix + "*",
		rolePermissionsKeyPrefix + "*",
		userPermissionsKeyPrefix + "*",
		permissionCheckKeyPrefix + "*",
	}
	
	totalDeleted := 0
	
	for _, pattern := range patterns {
		var cursor uint64
		var keys []string
		
		// 使用 SCAN 查找所有匹配的 key
		for {
			var batch []string
			var err error
			batch, cursor, err = c.redis.Scan(ctx, cursor, pattern, 1000).Result()
			if err != nil {
				c.logger.Error("Failed to scan keys",
					zap.String("pattern", pattern),
					zap.Error(err))
				continue
			}
			
			keys = append(keys, batch...)
			
			if cursor == 0 {
				break
			}
		}
		
		// 批量删除
		if len(keys) > 0 {
			if err := c.redis.Del(ctx, keys...).Err(); err != nil {
				c.logger.Error("Failed to delete keys",
					zap.String("pattern", pattern),
					zap.Error(err))
			} else {
				totalDeleted += len(keys)
				c.logger.Debug("Keys deleted",
					zap.String("pattern", pattern),
					zap.Int("count", len(keys)))
			}
		}
	}
	
	c.logger.Warn("All RBAC cache invalidated",
		zap.Int("total_deleted", totalDeleted))
	
	return nil
}

// === 统计信息 ===

// GetCacheStats 获取缓存统计信息
func (c *RBACCache) GetCacheStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)
	
	patterns := map[string]string{
		"user_roles":        userRolesKeyPrefix + "*",
		"role_permissions":  rolePermissionsKeyPrefix + "*",
		"user_permissions":  userPermissionsKeyPrefix + "*",
		"permission_checks": permissionCheckKeyPrefix + "*",
	}
	
	for name, pattern := range patterns {
		count := int64(0)
		var cursor uint64
		
		for {
			keys, newCursor, err := c.redis.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				c.logger.Error("Failed to scan keys",
					zap.String("pattern", pattern),
					zap.Error(err))
				break
			}
			
			count += int64(len(keys))
			cursor = newCursor
			
			if cursor == 0 {
				break
			}
		}
		
		stats[name] = count
	}
	
	return stats, nil
}

