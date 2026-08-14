package Usecases_test

import (
	"context"
	"errors"
	"testing"

	"task_manager/Domain/dto"
	"task_manager/Domain/models"
	"task_manager/Mocks"
	"task_manager/Usecases"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskServiceSuite struct {
	suite.Suite
	mockRepo *Mocks.MockTaskRepository
	service  *Usecases.TaskService
	ownerID  bson.ObjectID
}

func (s *TaskServiceSuite) SetupTest() {
	s.mockRepo = new(Mocks.MockTaskRepository)
	s.service = Usecases.NewTaskService(s.mockRepo)
	s.ownerID = bson.NewObjectID()
}

func TestTaskServiceSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceSuite))
}

func (s *TaskServiceSuite) TestAddTaskService_Success() {
	s.mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(task models.Task) bool {
		return task.Title == "fake" &&
			task.Description == "fake for testing" &&
			task.DueDate == "take date" &&
			task.Status == "PENDING" &&
			task.Owner == s.ownerID
	})).Return(true, nil)

	req := dto.TaskRequest{
		Title:       "fake",
		Description: "fake for testing",
		DueDate:     "take date",
	}

	result, err := s.service.AddTask(context.Background(), req, s.ownerID.Hex())

	s.NoError(err)
	s.Equal("fake", result.Title)
	s.Equal("fake for testing", result.Description)
	s.Equal("PENDING", result.Status)
	s.NotEmpty(result.ID)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestAddTaskService_RepositoryFailure() {
	s.mockRepo.On("Create", mock.Anything, mock.Anything).Return(false, errors.New("db error"))

	req := dto.TaskRequest{Title: "fake"}

	result, err := s.service.AddTask(context.Background(), req, bson.NewObjectID().Hex())

	s.Error(err)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestAddTaskService_InvalidOwnerID() {
	req := dto.TaskRequest{Title: "fake"}

	result, err := s.service.AddTask(context.Background(), req, "invalid-hex")

	s.Error(err)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

func (s *TaskServiceSuite) TestListTaskService_Successes() {
	fakeTasks := []models.Task{
		{ID: bson.NewObjectID(), Title: "test1", Description: "description 1", DueDate: "2026-11-11", Status: "PENDING", Owner: s.ownerID},
		{ID: bson.NewObjectID(), Title: "test2", Description: "description 2", DueDate: "2026-11-12", Status: "COMPLETED", Owner: s.ownerID},
		{ID: bson.NewObjectID(), Title: "test3", Description: "description 3", DueDate: "2026-11-13", Status: "COMPLETED", Owner: s.ownerID},
	}

	s.mockRepo.On("List", mock.Anything, s.ownerID).Return(fakeTasks, nil)

	result, err := s.service.ListTasks(context.Background(), s.ownerID.Hex())

	s.Require().NoError(err)
	s.Require().Len(result, 3)

	s.Equal(fakeTasks[0].ID.Hex(), result[0].ID)
	s.Equal(fakeTasks[0].Title, result[0].Title)
	s.Equal(fakeTasks[1].Status, result[1].Status)
	s.Equal(fakeTasks[2].Title, result[2].Title)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestListTaskService_Failure() {
	s.mockRepo.On("List", mock.Anything, s.ownerID).Return([]models.Task{}, errors.New("not found"))

	result, err := s.service.ListTasks(context.Background(), s.ownerID.Hex())

	s.Require().Error(err)
	s.Require().Len(result, 0)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestListTaskService_InvalidOwnerID() {
	result, err := s.service.ListTasks(context.Background(), "invalid-hex")

	s.Error(err)
	s.Nil(result)

	s.mockRepo.AssertNotCalled(s.T(), "List", mock.Anything, mock.Anything)
}

func (s *TaskServiceSuite) TestGetTaskService_Success() {
	task := models.Task{
		ID:          bson.NewObjectID(),
		Title:       "task 1",
		Description: "desc",
		DueDate:     "2026-11-11",
		Status:      "PENDING",
		Owner:       s.ownerID,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil)

	result, err := s.service.GetTask(context.Background(), task.ID.Hex(), s.ownerID.Hex())

	s.Require().NoError(err)
	s.Equal(task.ID.Hex(), result.ID)
	s.Equal(task.Title, result.Title)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestGetTaskService_Unauthorized() {
	otherOwner := bson.NewObjectID()
	task := models.Task{
		ID:     bson.NewObjectID(),
		Title:  "private",
		Status: "PENDING",
		Owner:  otherOwner,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil)

	result, err := s.service.GetTask(context.Background(), task.ID.Hex(), s.ownerID.Hex())

	s.ErrorIs(err, Usecases.ErrUnauthorized)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestGetTaskService_RepositoryFailure() {
	taskID := bson.NewObjectID()

	s.mockRepo.On("Get", mock.Anything, taskID).Return(models.Task{}, errors.New("not found"))

	result, err := s.service.GetTask(context.Background(), taskID.Hex(), s.ownerID.Hex())

	s.Error(err)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestGetTaskService_InvalidID() {
	result, err := s.service.GetTask(context.Background(), "invalid-hex", s.ownerID.Hex())

	s.Error(err)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertNotCalled(s.T(), "Get", mock.Anything, mock.Anything)
}

func (s *TaskServiceSuite) TestUpdateTaskService_Success() {
	task := models.Task{
		ID:          bson.NewObjectID(),
		Title:       "old title",
		Description: "old desc",
		DueDate:     "2026-11-11",
		Status:      "PENDING",
		Owner:       s.ownerID,
	}

	updated := models.Task{
		ID:          task.ID,
		Title:       "new title",
		Description: "new desc",
		DueDate:     "2026-12-12",
		Status:      "COMPLETED",
		Owner:       s.ownerID,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil).Once()
	s.mockRepo.On("Update", mock.Anything, task.ID, mock.Anything).Return(true, nil).Once()
	s.mockRepo.On("Get", mock.Anything, task.ID).Return(updated, nil).Once()

	result, err := s.service.UpdateTask(context.Background(), task.ID.Hex(), s.ownerID.Hex(), dto.TaskUpdateRequest{
		Title:       "new title",
		Description: "new desc",
		DueDate:     "2026-12-12",
		Status:      "COMPLETED",
	})

	s.Require().NoError(err)
	s.Equal("new title", result.Title)
	s.Equal("COMPLETED", result.Status)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestUpdateTaskService_Unauthorized() {
	otherOwner := bson.NewObjectID()
	task := models.Task{
		ID:     bson.NewObjectID(),
		Title:  "private",
		Status: "PENDING",
		Owner:  otherOwner,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil)

	result, err := s.service.UpdateTask(context.Background(), task.ID.Hex(), s.ownerID.Hex(), dto.TaskUpdateRequest{Title: "new"})

	s.ErrorIs(err, Usecases.ErrUnauthorized)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestUpdateTaskService_NoFieldsToUpdate() {
	task := models.Task{
		ID:     bson.NewObjectID(),
		Title:  "title",
		Status: "PENDING",
		Owner:  s.ownerID,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil)

	result, err := s.service.UpdateTask(context.Background(), task.ID.Hex(), s.ownerID.Hex(), dto.TaskUpdateRequest{})

	s.Error(err)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
}

func (s *TaskServiceSuite) TestUpdateTaskService_TaskNotFound() {
	task := models.Task{
		ID:     bson.NewObjectID(),
		Title:  "title",
		Status: "PENDING",
		Owner:  s.ownerID,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil)
	s.mockRepo.On("Update", mock.Anything, task.ID, mock.Anything).Return(false, nil)

	result, err := s.service.UpdateTask(context.Background(), task.ID.Hex(), s.ownerID.Hex(), dto.TaskUpdateRequest{Title: "new"})

	s.Error(err)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestUpdateTaskService_RepositoryFailure() {
	taskID := bson.NewObjectID()

	s.mockRepo.On("Get", mock.Anything, taskID).Return(models.Task{}, errors.New("not found"))

	result, err := s.service.UpdateTask(context.Background(), taskID.Hex(), s.ownerID.Hex(), dto.TaskUpdateRequest{Title: "new"})

	s.Error(err)
	s.Equal(dto.TaskResponse{}, result)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestDeleteTaskService_Success() {
	task := models.Task{
		ID:     bson.NewObjectID(),
		Title:  "title",
		Status: "PENDING",
		Owner:  s.ownerID,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil)
	s.mockRepo.On("Delete", mock.Anything, task.ID).Return(true, nil)

	s.NoError(s.service.DeleteTask(context.Background(), task.ID.Hex(), s.ownerID.Hex()))

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestDeleteTaskService_Unauthorized() {
	otherOwner := bson.NewObjectID()
	task := models.Task{
		ID:     bson.NewObjectID(),
		Title:  "private",
		Status: "PENDING",
		Owner:  otherOwner,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil)

	err := s.service.DeleteTask(context.Background(), task.ID.Hex(), s.ownerID.Hex())

	s.ErrorIs(err, Usecases.ErrUnauthorized)

	s.mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
}

func (s *TaskServiceSuite) TestDeleteTaskService_TaskNotFound() {
	task := models.Task{
		ID:     bson.NewObjectID(),
		Title:  "title",
		Status: "PENDING",
		Owner:  s.ownerID,
	}

	s.mockRepo.On("Get", mock.Anything, task.ID).Return(task, nil)
	s.mockRepo.On("Delete", mock.Anything, task.ID).Return(false, nil)

	err := s.service.DeleteTask(context.Background(), task.ID.Hex(), s.ownerID.Hex())

	s.Error(err)

	s.mockRepo.AssertExpectations(s.T())
}

func (s *TaskServiceSuite) TestDeleteTaskService_RepositoryFailure() {
	taskID := bson.NewObjectID()

	s.mockRepo.On("Get", mock.Anything, taskID).Return(models.Task{}, errors.New("not found"))

	err := s.service.DeleteTask(context.Background(), taskID.Hex(), s.ownerID.Hex())

	s.Error(err)

	s.mockRepo.AssertExpectations(s.T())
}
