package Usecases_test

import (
	"context"
	"errors"
	"testing"

	"task_manager/Domain/dto"
	"task_manager/Domain/models"
	"task_manager/Infrastructure/utils"
	"task_manager/Mocks"
	"task_manager/Usecases"
	"task_manager/repository"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserServiceSuite struct {
	suite.Suite
	mockRepo *Mocks.MockUserRepository
	service  *Usecases.UserService
}

func (s *UserServiceSuite) SetupTest() {
	s.T().Setenv("SECRET_KEY", "test-secret")
	s.mockRepo = new(Mocks.MockUserRepository)
	s.service = Usecases.NewUserService(s.mockRepo)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) TestRegisterService_Success() {
	s.mockRepo.On("FindByEmail", mock.Anything, "nafiad@example.com").Return(models.User{}, repository.ErrNotFound)
	s.mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(user models.User) bool {
		return user.Name == "Nafiad" &&
			user.Email == "nafiad@example.com" &&
			user.Role == "user" &&
			user.PasswordHash != "" &&
			user.PasswordHash != "password123"
	})).Return(true, nil)

	req := dto.RegisterRequest{
		Name:     "Nafiad",
		Email:    "nafiad@example.com",
		Password: "password123",
	}

	result, err := s.service.Register(context.Background(), req)

	s.NoError(err)
	s.NotEmpty(result.ID)
	s.Equal("Nafiad", result.Name)
	s.Equal("nafiad@example.com", result.Email)
	s.Equal("user", result.Role)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestRegisterService_EmailAlreadyTaken() {
	s.mockRepo.On("FindByEmail", mock.Anything, "nafiad@example.com").Return(models.User{
		ID:    bson.NewObjectID(),
		Email: "nafiad@example.com",
	}, nil)

	req := dto.RegisterRequest{
		Name:     "Nafiad",
		Email:    "nafiad@example.com",
		Password: "password123",
	}

	result, err := s.service.Register(context.Background(), req)

	s.Error(err)
	s.Equal("Email already taken.", err.Error())
	s.Equal(dto.UserResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestRegisterService_RepositoryFailure() {
	s.mockRepo.On("FindByEmail", mock.Anything, "nafiad@example.com").Return(models.User{}, errors.New("db error"))

	req := dto.RegisterRequest{
		Name:     "Nafiad",
		Email:    "nafiad@example.com",
		Password: "password123",
	}

	result, err := s.service.Register(context.Background(), req)

	s.Error(err)
	s.Equal("db error", err.Error())
	s.Equal(dto.UserResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestRegisterService_CreateFailure() {
	s.mockRepo.On("FindByEmail", mock.Anything, "nafiad@example.com").Return(models.User{}, repository.ErrNotFound)
	s.mockRepo.On("Create", mock.Anything, mock.Anything).Return(false, errors.New("db insert error"))

	req := dto.RegisterRequest{
		Name:     "Nafiad",
		Email:    "nafiad@example.com",
		Password: "password123",
	}

	result, err := s.service.Register(context.Background(), req)

	s.Error(err)
	s.Equal("db insert error", err.Error())
	s.Equal(dto.UserResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestLoginService_Success() {
	user := models.User{
		ID:           bson.NewObjectID(),
		Name:         "Nafiad",
		Email:        "nafiad@example.com",
		PasswordHash: hashPassword(s.T(), "password123"),
		Role:         "user",
	}

	s.mockRepo.On("FindByEmail", mock.Anything, "nafiad@example.com").Return(user, nil)

	token, err := s.service.Login(context.Background(), dto.LoginRequest{
		Email:    "nafiad@example.com",
		Password: "password123",
	})

	s.NoError(err)
	s.NotEmpty(token)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestLoginService_UserNotFound() {
	s.mockRepo.On("FindByEmail", mock.Anything, "unknown@example.com").Return(models.User{}, repository.ErrNotFound)

	token, err := s.service.Login(context.Background(), dto.LoginRequest{
		Email:    "unknown@example.com",
		Password: "password123",
	})

	s.Error(err)
	s.Equal("Not Found", err.Error())
	s.Empty(token)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestLoginService_InvalidPassword() {
	user := models.User{
		ID:           bson.NewObjectID(),
		Name:         "Nafiad",
		Email:        "nafiad@example.com",
		PasswordHash: hashPassword(s.T(), "password123"),
		Role:         "user",
	}

	s.mockRepo.On("FindByEmail", mock.Anything, "nafiad@example.com").Return(user, nil)

	token, err := s.service.Login(context.Background(), dto.LoginRequest{
		Email:    "nafiad@example.com",
		Password: "wrong-password",
	})

	s.Error(err)
	s.Equal("Invalid Credentials", err.Error())
	s.Empty(token)

	s.mockRepo.AssertExpectations(s.T())
}

func hashPassword(t *testing.T, password string) string {
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return hash
}
