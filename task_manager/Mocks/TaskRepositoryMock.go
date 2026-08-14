package Mocks

import (
	"context"
	"task_manager/Domain/models"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task models.Task) (bool, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockTaskRepository) List(ctx context.Context, ownerID bson.ObjectID) ([]models.Task, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) Get(ctx context.Context, objectID bson.ObjectID) (models.Task, error) {
	args := m.Called(ctx, objectID)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, objectID bson.ObjectID, set bson.M) (bool, error) {
	args := m.Called(ctx, objectID, set)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockTaskRepository) Delete(ctx context.Context, objectID bson.ObjectID) (bool, error) {
	args := m.Called(ctx, objectID)
	return args.Get(0).(bool), args.Error(1)
}
