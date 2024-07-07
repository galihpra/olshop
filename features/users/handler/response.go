package handler

type LoginResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email" form:"email"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

type UserResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email" form:"email"`
	Image    string `json:"image"`
}
