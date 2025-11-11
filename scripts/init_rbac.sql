-- RBAC 权限系统初始化脚本

-- 创建表
CREATE TABLE IF NOT EXISTS `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `display_name` varchar(100) NOT NULL COMMENT '显示名称',
  `description` varchar(500) DEFAULT NULL COMMENT '角色描述',
  `status` int NOT NULL DEFAULT '1' COMMENT '状态: 1-启用 0-禁用',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

CREATE TABLE IF NOT EXISTS `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(100) NOT NULL COMMENT '权限编码',
  `name` varchar(100) NOT NULL COMMENT '权限名称',
  `resource` varchar(50) NOT NULL COMMENT '资源',
  `action` varchar(50) NOT NULL COMMENT '操作',
  `description` varchar(500) DEFAULT NULL COMMENT '权限描述',
  `status` int NOT NULL DEFAULT '1' COMMENT '状态: 1-启用 0-禁用',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

CREATE TABLE IF NOT EXISTS `role_permissions` (
  `role_id` bigint unsigned NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `fk_role_permissions_permission` (`permission_id`),
  CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

CREATE TABLE IF NOT EXISTS `user_roles` (
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `fk_user_roles_role` (`role_id`),
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 插入默认权限
INSERT INTO `permissions` (`code`, `name`, `resource`, `action`, `description`, `status`, `created_at`) VALUES
-- 用户管理权限
('user:read', '查看用户', 'user', 'read', '查看用户列表和详情', 1, NOW()),
('user:write', '编辑用户', 'user', 'write', '编辑用户信息、状态、重置密码', 1, NOW()),
('user:delete', '删除用户', 'user', 'delete', '删除用户', 1, NOW()),

-- 角色权限管理
('rbac:manage', 'RBAC管理', 'rbac', 'manage', '管理角色和权限', 1, NOW()),

-- 统计信息权限
('statistics:read', '查看统计', 'statistics', 'read', '查看统计信息', 1, NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 插入默认角色
INSERT INTO `roles` (`name`, `display_name`, `description`, `status`, `created_at`) VALUES
('superadmin', '超级管理员', '拥有所有权限的超级管理员', 1, NOW()),
('admin', '管理员', '拥有大部分管理权限', 1, NOW()),
('editor', '编辑员', '可以查看和编辑用户，但不能删除', 1, NOW()),
('viewer', '查看者', '只能查看信息，不能编辑', 1, NOW())
ON DUPLICATE KEY UPDATE `display_name` = VALUES(`display_name`);

-- 为角色分配权限

-- 超级管理员：所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 
  (SELECT id FROM roles WHERE name = 'superadmin'),
  id,
  NOW()
FROM permissions
ON DUPLICATE KEY UPDATE `created_at` = VALUES(`created_at`);

-- 管理员：除 RBAC 管理外的所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 
  (SELECT id FROM roles WHERE name = 'admin'),
  id,
  NOW()
FROM permissions
WHERE code IN ('user:read', 'user:write', 'user:delete', 'statistics:read')
ON DUPLICATE KEY UPDATE `created_at` = VALUES(`created_at`);

-- 编辑员：查看和编辑用户
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 
  (SELECT id FROM roles WHERE name = 'editor'),
  id,
  NOW()
FROM permissions
WHERE code IN ('user:read', 'user:write', 'statistics:read')
ON DUPLICATE KEY UPDATE `created_at` = VALUES(`created_at`);

-- 查看者：只能查看
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 
  (SELECT id FROM roles WHERE name = 'viewer'),
  id,
  NOW()
FROM permissions
WHERE code IN ('user:read', 'statistics:read')
ON DUPLICATE KEY UPDATE `created_at` = VALUES(`created_at`);

-- 为初始管理员用户（ID=1）分配超级管理员角色
-- 注意：这假设你的管理员用户 ID 是 1，请根据实际情况调整
INSERT INTO `user_roles` (`user_id`, `role_id`, `created_at`)
SELECT 
  1,
  (SELECT id FROM roles WHERE name = 'superadmin'),
  NOW()
ON DUPLICATE KEY UPDATE `created_at` = VALUES(`created_at`);

-- 查看结果
SELECT '============ 角色列表 ============' as '';
SELECT id, name, display_name, description FROM roles;

SELECT '============ 权限列表 ============' as '';
SELECT id, code, name, resource, action FROM permissions;

SELECT '============ 角色权限关系 ============' as '';
SELECT 
  r.name as role_name,
  p.code as permission_code,
  p.name as permission_name
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
ORDER BY r.name, p.code;

SELECT '============ 用户角色关系 ============' as '';
SELECT 
  ur.user_id,
  u.username,
  r.name as role_name,
  r.display_name as role_display_name
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id;

