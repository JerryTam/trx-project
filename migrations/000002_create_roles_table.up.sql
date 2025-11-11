-- 创建角色表
CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '角色名称',
    `display_name` VARCHAR(100) NOT NULL COMMENT '显示名称',
    `description` VARCHAR(500) NULL COMMENT '角色描述',
    `status` INT NOT NULL DEFAULT 1 COMMENT '状态：1-启用 0-禁用',
    `created_at` DATETIME(3) NULL DEFAULT NULL,
    `updated_at` DATETIME(3) NULL DEFAULT NULL,
    `deleted_at` DATETIME(3) NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_roles_name` (`name`),
    INDEX `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

