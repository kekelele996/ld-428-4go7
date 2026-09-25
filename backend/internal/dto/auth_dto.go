package dto

// LoginRequest 登录。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=Admin Curator Artist Viewer"`
}

// LoginResponse 登录响应。
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserSummary `json:"user"`
}

// UserSummary 用户摘要。
type UserSummary struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}
