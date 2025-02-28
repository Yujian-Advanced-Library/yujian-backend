package model

type LoginRequestDTO struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	BaseResp
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type RegisterRequestDTO struct {
	UserName        string `json:"user_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type RegisterResponseDTO struct {
	BaseResp
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
