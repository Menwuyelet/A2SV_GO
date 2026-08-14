package Mocks

import (
	"context"
	"task_manager/Domain/dto"

	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(ctx context.Context, req dto.RegisterRequest) (dto.UserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}

func (m *MockUserUsecase) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(string), args.Error(1)
}
