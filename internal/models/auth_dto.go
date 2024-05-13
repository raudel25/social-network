package models

type RegisterRequest struct {
	Name     string     `json:"name" binding:"required"`
	Username string     `json:"username" binding:"required"`
	Email    string     `json:"email" binding:"required"`
	Password string     `json:"password" binding:"required"`
	Profile  ProfileDto `json:"profile" binding:"required"`
}

type LoginResponse struct {
	Username string  `json:"username" binding:"required"`
	Profile  Profile `json:"profile" binding:"required"`
	Token    string  `json:"token" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
