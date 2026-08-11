package dto

import "go.mongodb.org/mongo-driver/v2/bson"

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID    bson.ObjectID `json:"id"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
	Role  string        `json:"role"`
}
