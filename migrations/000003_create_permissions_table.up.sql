-- 创建权限表
CREATE TABLE IF NOT EXISTS `permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(100) NOT NULL COMMENT '权限编码',
    `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
    `resource` VARCHAR(50) NOT NULL COMMENT '资源',
    `action` VARCHAR(50) NOT NULL COMMENT '操作',
    `description` VARCHAR(500) NULL COMMENT '权限描述',
    `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-启用 0-禁用',
    `created_at` DATETIME(3) NULL DEFAULT NULL,
    `updated_at` DATETIME(3) NULL DEFAULT NULL,
    `deleted_at` DATETIME(3) NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_permissions_code` (`code`),
    INDEX `idx_permissions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

