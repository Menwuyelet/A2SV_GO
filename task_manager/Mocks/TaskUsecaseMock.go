package Mocks

import (
	"context"
	"task_manager/Domain/dto"

	"github.com/stretchr/testify/mock"
)

type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) AddTask(ctx context.Context, task dto.TaskRequest, owner string) (dto.TaskResponse, error) {
	args := m.Called(ctx, task, owner)
	return args.Get(0).(dto.TaskResponse), args.Error(1)
}

func (m *MockTaskUsecase) ListTasks(ctx context.Context, owner string) ([]dto.TaskResponse, error) {
	args := m.Called(ctx, owner)
	return args.Get(0).([]dto.TaskResponse), args.Error(1)
}

func (m *MockTaskUsecase) GetTask(ctx context.Context, id string, owner string) (dto.TaskResponse, error) {
	args := m.Called(ctx, id, owner)
	return args.Get(0).(dto.TaskResponse), args.Error(1)
}

func (m *MockTaskUsecase) UpdateTask(ctx context.Context, id string, owner string, updatedTask dto.TaskUpdateRequest) (dto.TaskResponse, error) {
	args := m.Called(ctx, id, owner, updatedTask)
	return args.Get(0).(dto.TaskResponse), args.Error(1)
}

func (m *MockTaskUsecase) DeleteTask(ctx context.Context, id string, owner string) error {
	args := m.Called(ctx, id, owner)
	return args.Error(0)
}
