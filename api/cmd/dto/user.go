package dto

type UserRegisterDTO struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Handler  string `json:"handler" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=0,max=60"`
}

type GetUsersByUsernameDTO struct {
	Handle string `query:"handle"`
}
