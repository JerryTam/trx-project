package dto

import (
	"time"
	"trx-project/internal/model"
)

// UserDTO 用户响应 DTO
// 用于 API 响应，不包含敏感信息（如密码）
type UserDTO struct {
	ID        uint      `json:"id"`         // 用户ID
	Username  string    `json:"username"`   // 用户名
	Email     string    `json:"email"`      // 邮箱
	Status    int       `json:"status"`     // 状态
	StatusText string   `json:"status_text"` // 状态文本
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
	// 注意：不包含 password 和 deleted_at 字段
}

// UserProfileDTO 用户个人信息 DTO
// 用于用户个人信息接口，包含更多详细信息
type UserProfileDTO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	StatusText string   `json:"status_text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 可以添加更多个人信息字段
	// Avatar    string    `json:"avatar"`    // 头像
	// Phone     string    `json:"phone"`     // 手机号
	// Nickname  string    `json:"nickname"`  // 昵称
}

// UserListDTO 用户列表 DTO（简化版）
// 用于列表接口，只包含必要字段
type UserListDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	StatusText string `json:"status_text"`
}

// ToUserDTO 将 Model 转换为 DTO
func ToUserDTO(user *model.User) *UserDTO {
	if user == nil {
		return nil
	}

	return &UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		StatusText: getUserStatusText(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserProfileDTO 将 Model 转换为个人信息 DTO
func ToUserProfileDTO(user *model.User) *UserProfileDTO {
	if user == nil {
		return nil
	}

	return &UserProfileDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		StatusText: getUserStatusText(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserListDTO 将 Model 转换为列表 DTO
func ToUserListDTO(user *model.User) *UserListDTO {
	if user == nil {
		return nil
	}

	return &UserListDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		StatusText: getUserStatusText(user.Status),
	}
}

// ToUserDTOList 批量转换 Model 列表为 DTO 列表
func ToUserDTOList(users []model.User) []*UserDTO {
	if len(users) == 0 {
		return []*UserDTO{}
	}

	result := make([]*UserDTO, 0, len(users))
	for i := range users {
		result = append(result, ToUserDTO(&users[i]))
	}
	return result
}

// ToUserListDTOList 批量转换 Model 列表为列表 DTO
func ToUserListDTOList(users []model.User) []*UserListDTO {
	if len(users) == 0 {
		return []*UserListDTO{}
	}

	result := make([]*UserListDTO, 0, len(users))
	for i := range users {
		result = append(result, ToUserListDTO(&users[i]))
	}
	return result
}

// getUserStatusText 获取用户状态文本（辅助函数）
func getUserStatusText(status int) string {
	// 根据实际 User 模型的状态定义
	// 1: 活跃, 0: 禁用
	statusTextMap := map[int]string{
		1: "正常",
		0: "禁用",
	}
	if text, ok := statusTextMap[status]; ok {
		return text
	}
	return "未知"
}

