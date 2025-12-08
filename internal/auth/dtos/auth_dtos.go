package dtos

type RegisterDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponseDto struct {
	Token string `json:"token"`
}
