package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    string
	Name  string
	Email string
	Role  string
}
