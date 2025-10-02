package dto

type UserRequest struct {
	Username string `json:"username" validate:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
