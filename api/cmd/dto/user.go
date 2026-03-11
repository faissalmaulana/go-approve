package dto

type UserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Handler  string `json:"handler"`
	Password string `json:"password"`
}
