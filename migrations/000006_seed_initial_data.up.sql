-- 初始化基础权限数据
INSERT INTO `permissions` (`code`, `name`, `resource`, `action`, `description`, `status`, `created_at`, `updated_at`) VALUES
-- 用户管理权限
('user:read', '查看用户', 'user', 'read', '查看用户信息和列表', 1, NOW(), NOW()),
('user:write', '编辑用户', 'user', 'write', '创建和编辑用户信息', 1, NOW(), NOW()),
('user:delete', '删除用户', 'user', 'delete', '删除用户', 1, NOW(), NOW()),

-- RBAC 管理权限
('rbac:manage', 'RBAC管理', 'rbac', 'manage', '管理角色和权限', 1, NOW(), NOW()),

-- 统计信息权限
('statistics:read', '查看统计', 'statistics', 'read', '查看统计信息', 1, NOW(), NOW());

-- 初始化基础角色数据
INSERT INTO `roles` (`name`, `display_name`, `description`, `status`, `created_at`, `updated_at`) VALUES
('superadmin', '超级管理员', '拥有所有权限的超级管理员', 1, NOW(), NOW()),
('admin', '管理员', '拥有大部分管理权限', 1, NOW(), NOW()),
('editor', '编辑', '拥有基本编辑权限', 1, NOW(), NOW()),
('viewer', '查看者', '只有查看权限', 1, NOW(), NOW());

-- 为超级管理员分配所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 
    (SELECT `id` FROM `roles` WHERE `name` = 'superadmin'),
    `id`,
    NOW()
FROM `permissions`;

-- 为管理员分配权限（除了 rbac:manage）
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 
    (SELECT `id` FROM `roles` WHERE `name` = 'admin'),
    `id`,
    NOW()
FROM `permissions`
WHERE `code` IN ('user:read', 'user:write', 'user:delete', 'statistics:read');

-- 为编辑分配权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 
    (SELECT `id` FROM `roles` WHERE `name` = 'editor'),
    `id`,
    NOW()
FROM `permissions`
WHERE `code` IN ('user:read', 'user:write');

-- 为查看者分配权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 
    (SELECT `id` FROM `roles` WHERE `name` = 'viewer'),
    `id`,
    NOW()
FROM `permissions`
WHERE `code` = 'user:read';

