package main

import (
	"fmt"
	"time"
	"trx-project/pkg/jwt"
)

func main() {
	// JWT 配置（需要与 config/config.yaml 中的配置一致）
	config := jwt.Config{
		Secret:     "your-secret-key-change-in-production",
		Issuer:     "trx-project",
		ExpireTime: 24 * time.Hour, // 管理员 Token 1天有效期
	}

	// 生成超级管理员 Token
	token, err := jwt.GenerateToken(1, "admin", "superadmin", config)
	if err != nil {
		fmt.Printf("❌ Failed to generate token: %v\n", err)
		return
	}

	fmt.Println("✅ 管理员 Token 生成成功!")
	fmt.Println()
	fmt.Println("Token 信息:")
	fmt.Println("  User ID:  1")
	fmt.Println("  Username: admin")
	fmt.Println("  Role:     superadmin")
	fmt.Println("  Expires:  1 day")
	fmt.Println()
	fmt.Println("Token:")
	fmt.Println(token)
	fmt.Println()
	fmt.Println("使用示例:")
	fmt.Printf("  export ADMIN_TOKEN=\"%s\"\n", token)
	fmt.Println("  curl -H \"Authorization: Bearer $ADMIN_TOKEN\" http://localhost:8081/api/v1/admin/users")
}
