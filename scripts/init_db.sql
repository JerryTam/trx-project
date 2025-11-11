-- 创建数据库
CREATE DATABASE IF NOT EXISTS trx_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE trx_db;

-- 用户表将通过 GORM 自动迁移创建
-- 此脚本仅作为参考

-- 示例：手动创建用户表（可选，GORM 会自动创建）
-- CREATE TABLE IF NOT EXISTS users (
--     id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
--     created_at DATETIME(3) NULL,
--     updated_at DATETIME(3) NULL,
--     deleted_at DATETIME(3) NULL,
--     username VARCHAR(50) NOT NULL,
--     email VARCHAR(100) NOT NULL,
--     password VARCHAR(255) NOT NULL,
--     status INT DEFAULT 1,
--     UNIQUE KEY idx_username (username),
--     UNIQUE KEY idx_email (email),
--     KEY idx_deleted_at (deleted_at)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

