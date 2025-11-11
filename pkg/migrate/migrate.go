package migrate

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

// Config 迁移配置
type Config struct {
	MigrationsPath string // 迁移文件路径，例如 "file://migrations"
	DatabaseURL    string // 数据库连接字符串
	Logger         *zap.Logger
}

// Migrator 数据库迁移管理器
type Migrator struct {
	migrate *migrate.Migrate
	logger  *zap.Logger
}

// NewMigrator 创建新的迁移管理器
func NewMigrator(cfg *Config) (*Migrator, error) {
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}

	// 打开数据库连接
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 创建 MySQL 驱动
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create mysql driver: %w", err)
	}

	// 创建迁移实例
	m, err := migrate.NewWithDatabaseInstance(
		cfg.MigrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &Migrator{
		migrate: m,
		logger:  cfg.Logger,
	}, nil
}

// Up 执行所有待执行的迁移（升级到最新版本）
func (m *Migrator) Up() error {
	m.logger.Info("Running database migrations...")

	err := m.migrate.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		m.logger.Info("No migrations to run (already up to date)")
		return nil
	}

	version, dirty, err := m.migrate.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		m.logger.Warn("Failed to get migration version", zap.Error(err))
	} else {
		m.logger.Info("Database migrations completed successfully",
			zap.Uint("version", version),
			zap.Bool("dirty", dirty))
	}

	return nil
}

// Down 回滚一个迁移版本
func (m *Migrator) Down() error {
	m.logger.Info("Rolling back database migration...")

	err := m.migrate.Steps(-1)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	version, dirty, _ := m.migrate.Version()
	m.logger.Info("Database migration rolled back successfully",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty))

	return nil
}

// Migrate 迁移到指定版本
func (m *Migrator) Migrate(version uint) error {
	m.logger.Info("Migrating to version", zap.Uint("target_version", version))

	err := m.migrate.Migrate(version)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		m.logger.Info("Already at target version", zap.Uint("version", version))
		return nil
	}

	m.logger.Info("Migration to version completed", zap.Uint("version", version))
	return nil
}

// Version 获取当前迁移版本
func (m *Migrator) Version() (uint, bool, error) {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			m.logger.Info("No migrations have been run yet")
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// Force 强制设置迁移版本（用于修复脏状态）
func (m *Migrator) Force(version int) error {
	m.logger.Warn("Forcing migration version", zap.Int("version", version))

	err := m.migrate.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	m.logger.Info("Migration version forced successfully", zap.Int("version", version))
	return nil
}

// Drop 删除所有表（危险操作！）
func (m *Migrator) Drop() error {
	m.logger.Warn("Dropping all database tables...")

	err := m.migrate.Drop()
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	m.logger.Info("All database tables dropped successfully")
	return nil
}

// Close 关闭迁移器
func (m *Migrator) Close() error {
	srcErr, dbErr := m.migrate.Close()
	if srcErr != nil {
		return fmt.Errorf("failed to close source: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close database: %w", dbErr)
	}
	return nil
}

