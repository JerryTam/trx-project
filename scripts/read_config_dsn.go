// +build dev

package main

import (
	"fmt"
	"os"
	"trx-project/pkg/config"
)

// 读取配置文件并生成数据库 DSN
// 用法: go run -tags=dev scripts/read_config_dsn.go <config_file>
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "用法: %s <config_file>\n", os.Args[0])
		os.Exit(1)
	}

	cfg, err := config.Load(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取配置失败: %v\n", err)
		os.Exit(1)
	}

	// 生成 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.MySQL.Username,
		cfg.Database.MySQL.Password,
		cfg.Database.MySQL.Host,
		cfg.Database.MySQL.Port,
		cfg.Database.MySQL.Database,
	)

	fmt.Print(dsn)
}

