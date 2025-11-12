package main

import (
	"fmt"
	"time"
	"trx-project/pkg/jwt"
)

func main() {
	config := jwt.Config{
		Secret:     "your-secret-key-change-in-production",
		Issuer:     "trx-project",
		ExpireTime: 24 * time.Hour,
	}

	// 生成超级管理员 Token（用户ID=1，角色=superadmin）
	fmt.Println("==================== 超级管理员 Token ====================")
	token, _ := jwt.GenerateToken(1, "admin", "superadmin", config)
	fmt.Println("✅ 超级管理员 Token 生成成功!")
	fmt.Println()
	fmt.Println("Token 信息:")
	fmt.Println("  User ID:  1")
	fmt.Println("  Username: admin")
	fmt.Println("  Role:     superadmin")
	fmt.Println("  权限:     所有权限（包括 RBAC 管理）")
	fmt.Println("  Expires:  1 day")
	fmt.Println()
	fmt.Println("Token:")
	fmt.Println(token)
	fmt.Println()
	fmt.Println("使用示例:")
	fmt.Printf("  export ADMIN_TOKEN=\"%s\"\n", token)
	fmt.Println("  curl -H \"Authorization: Bearer $ADMIN_TOKEN\" http://localhost:8081/api/v1/admin/users")
	fmt.Println()

	// 生成普通管理员 Token（用户ID=2，角色=admin）
	fmt.Println("==================== 普通管理员 Token ====================")
	token2, _ := jwt.GenerateToken(2, "manager", "admin", config)
	fmt.Println("✅ 普通管理员 Token 生成成功!")
	fmt.Println()
	fmt.Println("Token 信息:")
	fmt.Println("  User ID:  2")
	fmt.Println("  Username: manager")
	fmt.Println("  Role:     admin")
	fmt.Println("  权限:     user:read, user:write, user:delete, statistics:read")
	fmt.Println("  Expires:  1 day")
	fmt.Println()
	fmt.Println("Token:")
	fmt.Println(token2)
	fmt.Println()

	// 生成编辑员 Token（用户ID=3，角色=editor）
	fmt.Println("==================== 编辑员 Token ====================")
	token3, _ := jwt.GenerateToken(3, "editor", "admin", config)
	fmt.Println("✅ 编辑员 Token 生成成功!")
	fmt.Println()
	fmt.Println("Token 信息:")
	fmt.Println("  User ID:  3")
	fmt.Println("  Username: editor")
	fmt.Println("  Role:     admin (需要通过数据库分配 editor 角色)")
	fmt.Println("  权限:     user:read, user:write, statistics:read")
	fmt.Println("  Expires:  1 day")
	fmt.Println()
	fmt.Println("Token:")
	fmt.Println(token3)
	fmt.Println()

	fmt.Println("==================== 注意事项 ====================")
	fmt.Println("1. 请先运行 scripts/init_rbac.sql 初始化 RBAC 数据")
	fmt.Println("2. Token 中的 role 只用于基本认证，实际权限由数据库中的角色权限决定")
	fmt.Println("3. 超级管理员（ID=1）已自动分配 superadmin 角色")
	fmt.Println("4. 其他用户需要通过 API 或数据库手动分配角色")
	fmt.Println()
	fmt.Println("==================== 快速测试 ====================")
	fmt.Println("# 1. 初始化 RBAC 数据")
	fmt.Println("mysql -u root -p trx_dev < scripts/init_rbac.sql")
	fmt.Println()
	fmt.Println("# 2. 使用超级管理员 Token 测试")
	fmt.Println("export ADMIN_TOKEN=\"" + token + "\"")
	fmt.Println()
	fmt.Println("# 3. 查看角色列表")
	fmt.Println("curl -H \"Authorization: Bearer $ADMIN_TOKEN\" http://localhost:8081/api/v1/admin/rbac/roles")
	fmt.Println()
	fmt.Println("# 4. 查看权限列表")
	fmt.Println("curl -H \"Authorization: Bearer $ADMIN_TOKEN\" http://localhost:8081/api/v1/admin/rbac/permissions")
	fmt.Println()
	fmt.Println("# 5. 查看用户角色")
	fmt.Println("curl -H \"Authorization: Bearer $ADMIN_TOKEN\" http://localhost:8081/api/v1/admin/users/1/roles")
	fmt.Println()
	fmt.Println("# 6. 查看用户权限")
	fmt.Println("curl -H \"Authorization: Bearer $ADMIN_TOKEN\" http://localhost:8081/api/v1/admin/users/1/permissions")
}
