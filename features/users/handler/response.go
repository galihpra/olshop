package handler

type LoginResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email" form:"email"`
	Token    string `json:"token"`
}
