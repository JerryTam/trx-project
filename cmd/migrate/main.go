package main

import (
	"flag"
	"fmt"
	"os"

	"trx-project/pkg/config"
	"trx-project/pkg/logger"
	"trx-project/pkg/migrate"

	"go.uber.org/zap"
)

func main() {
	// 命令行参数
	configPath := flag.String("config", "config/config.yaml", "配置文件路径")
	command := flag.String("cmd", "", "迁移命令: up, down, version, force, drop, goto")
	version := flag.Int("version", 0, "目标版本号（用于 goto 和 force 命令）")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("❌ 加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.InitLogger(&cfg.Logger); err != nil {
		fmt.Printf("❌ 初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	log := logger.Logger

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		cfg.Database.MySQL.Username,
		cfg.Database.MySQL.Password,
		cfg.Database.MySQL.Host,
		cfg.Database.MySQL.Port,
		cfg.Database.MySQL.Database,
	)

	// 创建迁移器
	migrator, err := migrate.NewMigrator(&migrate.Config{
		MigrationsPath: "file://migrations",
		DatabaseURL:    dsn,
		Logger:         log,
	})
	if err != nil {
		log.Fatal("创建迁移器失败", zap.Error(err))
	}
	defer migrator.Close()

	// 执行命令
	switch *command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatal("执行迁移失败", zap.Error(err))
		}
		fmt.Println("✅ 数据库迁移完成")

	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatal("回滚迁移失败", zap.Error(err))
		}
		fmt.Println("✅ 迁移回滚完成")

	case "version":
		v, dirty, err := migrator.Version()
		if err != nil {
			log.Fatal("获取版本失败", zap.Error(err))
		}
		fmt.Printf("📌 当前迁移版本: %d (dirty: %v)\n", v, dirty)

	case "force":
		if *version == 0 {
			log.Fatal("force 命令需要指定版本号", zap.String("usage", "-version <version>"))
		}
		if err := migrator.Force(*version); err != nil {
			log.Fatal("强制设置版本失败", zap.Error(err))
		}
		fmt.Printf("✅ 强制设置版本为: %d\n", *version)

	case "goto":
		if *version == 0 {
			log.Fatal("goto 命令需要指定版本号", zap.String("usage", "-version <version>"))
		}
		if err := migrator.Migrate(uint(*version)); err != nil {
			log.Fatal("迁移到指定版本失败", zap.Error(err))
		}
		fmt.Printf("✅ 迁移到版本: %d\n", *version)

	case "drop":
		// 二次确认
		fmt.Print("⚠️  警告：此操作将删除所有表！请输入 'YES' 确认: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "YES" {
			fmt.Println("❌ 操作已取消")
			os.Exit(0)
		}
		if err := migrator.Drop(); err != nil {
			log.Fatal("删除表失败", zap.Error(err))
		}
		fmt.Println("✅ 所有表已删除")

	default:
		fmt.Println("❌ 未知命令:", *command)
		fmt.Println("\n可用命令:")
		fmt.Println("  up       - 执行所有待执行的迁移")
		fmt.Println("  down     - 回滚一个迁移版本")
		fmt.Println("  version  - 查看当前迁移版本")
		fmt.Println("  force    - 强制设置迁移版本 (需要 -version 参数)")
		fmt.Println("  goto     - 迁移到指定版本 (需要 -version 参数)")
		fmt.Println("  drop     - 删除所有表 (危险操作)")
		fmt.Println("\n示例:")
		fmt.Println("  go run cmd/migrate/main.go -cmd up")
		fmt.Println("  go run cmd/migrate/main.go -cmd version")
		fmt.Println("  go run cmd/migrate/main.go -cmd force -version 1")
		os.Exit(1)
	}
}

