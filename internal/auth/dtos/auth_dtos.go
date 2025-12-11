package dtos

type RegisterDto struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponseDto struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}
