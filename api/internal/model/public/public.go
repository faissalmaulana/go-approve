package public

// THIS IS modified field of the models that safe
// to seen by client

type UserPublic struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Handler  string `json:"handler"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
