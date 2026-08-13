package Usecases

import (
	"context"
	"errors"
	"task_manager/Domain/dto"
	"task_manager/Domain/models"
	"task_manager/Infrastructure/utils"
	"task_manager/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
}

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) Register(ctx context.Context, req dto.RegisterRequest) (dto.UserResponse, error) {
	_, err := us.repo.FindByEmail(ctx, req.Email)

	if err == nil {
		return dto.UserResponse{}, errors.New("Email already taken.")
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return dto.UserResponse{}, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return dto.UserResponse{}, err
	}

	userID := bson.NewObjectID()

	user := models.User{
		ID:           userID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	_, err = us.repo.Create(ctx, user)

	if err != nil {
		return dto.UserResponse{}, err
	}

	Newuser := dto.UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	return Newuser, nil
}

func (us *UserService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := us.repo.FindByEmail(ctx, req.Email)

	if err != nil {
		return "", errors.New("Not Found")
	}

	valid := utils.CheckPassword(user.PasswordHash, req.Password)

	if !valid {
		return "", errors.New("Invalid Credentials")
	}

	var AuthToken string

	AuthToken, err = utils.GenerateToken(user.ID.Hex(), user.Role)

	if err != nil {
		return "", err
	}

	return AuthToken, nil

}
