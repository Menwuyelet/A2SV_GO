package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"task_manager/Delivery/controllers"
	"task_manager/Domain/dto"
	"task_manager/Mocks"
	"task_manager/Usecases"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskControllerSuite struct {
	suite.Suite
	mockUsecase *Mocks.MockTaskUsecase
	controller  *controllers.TaskController
	ownerID     string
	taskID      string
}

func (s *TaskControllerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockUsecase = new(Mocks.MockTaskUsecase)
	s.controller = controllers.NewTaskController(s.mockUsecase)
	s.ownerID = "64c3f5e6c9e13a2f3b30a001"
	s.taskID = "64c3f5e6c9e13a2f3b30a002"
}

func TestTaskControllerSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerSuite))
}

func (s *TaskControllerSuite) newContext(method, body string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(method, "/", strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx, recorder
}

func (s *TaskControllerSuite) TestAddTask_Success() {
	req := dto.TaskRequest{Title: "task 1", Description: "desc", DueDate: "2026-11-11"}
	resp := dto.TaskResponse{ID: s.taskID, Title: "task 1", Description: "desc", DueDate: "2026-11-11", Status: "PENDING"}

	s.mockUsecase.On("AddTask", mock.Anything, req, s.ownerID).Return(resp, nil)

	ctx, recorder := s.newContext(http.MethodPost, `{"title":"task 1","description":"desc","duedate":"2026-11-11"}`)
	ctx.Set("userID", s.ownerID)

	s.controller.AddTask(ctx)

	s.Equal(http.StatusCreated, recorder.Code)
	var got dto.TaskResponse
	s.NoError(json.Unmarshal(recorder.Body.Bytes(), &got))
	s.Equal(resp, got)

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestAddTask_InvalidJSON() {
	ctx, recorder := s.newContext(http.MethodPost, `{"title":`)
	ctx.Set("userID", s.ownerID)

	s.controller.AddTask(ctx)

	s.Equal(http.StatusBadRequest, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "AddTask", mock.Anything, mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestAddTask_Unauthenticated() {
	ctx, recorder := s.newContext(http.MethodPost, `{"title":"task 1"}`)

	s.controller.AddTask(ctx)

	s.Equal(http.StatusUnauthorized, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "AddTask", mock.Anything, mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestAddTask_ServiceFailure() {
	req := dto.TaskRequest{Title: "task 1"}

	s.mockUsecase.On("AddTask", mock.Anything, req, s.ownerID).Return(dto.TaskResponse{}, errors.New("create failed"))

	ctx, recorder := s.newContext(http.MethodPost, `{"title":"task 1"}`)
	ctx.Set("userID", s.ownerID)

	s.controller.AddTask(ctx)

	s.Equal(http.StatusInternalServerError, recorder.Code)

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestListAllTasks_Success() {
	tasks := []dto.TaskResponse{
		{ID: s.taskID, Title: "task 1", Status: "PENDING"},
	}

	s.mockUsecase.On("ListTasks", mock.Anything, s.ownerID).Return(tasks, nil)

	ctx, recorder := s.newContext(http.MethodGet, "")
	ctx.Set("userID", s.ownerID)

	s.controller.ListAllTasks(ctx)

	s.Equal(http.StatusOK, recorder.Code)
	s.Contains(recorder.Body.String(), "task 1")

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestListAllTasks_Unauthenticated() {
	ctx, recorder := s.newContext(http.MethodGet, "")

	s.controller.ListAllTasks(ctx)

	s.Equal(http.StatusUnauthorized, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "ListTasks", mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestListAllTasks_ServiceFailure() {
	s.mockUsecase.On("ListTasks", mock.Anything, s.ownerID).Return([]dto.TaskResponse{}, errors.New("list failed"))

	ctx, recorder := s.newContext(http.MethodGet, "")
	ctx.Set("userID", s.ownerID)

	s.controller.ListAllTasks(ctx)

	s.Equal(http.StatusInternalServerError, recorder.Code)

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestRetrieveTask_Success() {
	task := dto.TaskResponse{ID: s.taskID, Title: "task 1", Status: "PENDING"}

	s.mockUsecase.On("GetTask", mock.Anything, s.taskID, s.ownerID).Return(task, nil)

	ctx, recorder := s.newContext(http.MethodGet, "")
	ctx.Set("userID", s.ownerID)
	ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

	s.controller.RetrieveTask(ctx)

	s.Equal(http.StatusOK, recorder.Code)
	s.Contains(recorder.Body.String(), "task 1")

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestRetrieveTask_Unauthenticated() {
	ctx, recorder := s.newContext(http.MethodGet, "")
	ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

	s.controller.RetrieveTask(ctx)

	s.Equal(http.StatusUnauthorized, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "GetTask", mock.Anything, mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestRetrieveTask_StatusMapping() {
	testCases := []struct {
		name       string
		mockErr    error
		wantStatus int
	}{
		{
			name:       "unauthorized owner",
			mockErr:    Usecases.ErrUnauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "task not found",
			mockErr:    errors.New("not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			mockUsecase := new(Mocks.MockTaskUsecase)
			controller := controllers.NewTaskController(mockUsecase)
			mockUsecase.On("GetTask", mock.Anything, s.taskID, s.ownerID).Return(dto.TaskResponse{}, tc.mockErr)

			ctx, recorder := s.newContext(http.MethodGet, "")
			ctx.Set("userID", s.ownerID)
			ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

			controller.RetrieveTask(ctx)

			s.Equal(tc.wantStatus, recorder.Code)
			mockUsecase.AssertExpectations(s.T())
		})
	}
}

func (s *TaskControllerSuite) TestUpdateTask_Success() {
	updateReq := dto.TaskUpdateRequest{Title: "new title"}
	task := dto.TaskResponse{ID: s.taskID, Title: "new title", Status: "PENDING"}

	s.mockUsecase.On("UpdateTask", mock.Anything, s.taskID, s.ownerID, updateReq).Return(task, nil)

	ctx, recorder := s.newContext(http.MethodPut, `{"title":"new title"}`)
	ctx.Set("userID", s.ownerID)
	ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

	s.controller.UpdateTask(ctx)

	s.Equal(http.StatusOK, recorder.Code)
	s.Contains(recorder.Body.String(), "new title")

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestUpdateTask_InvalidJSON() {
	ctx, recorder := s.newContext(http.MethodPut, `{"title":`)
	ctx.Set("userID", s.ownerID)
	ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

	s.controller.UpdateTask(ctx)

	s.Equal(http.StatusBadRequest, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "UpdateTask", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestUpdateTask_Unauthenticated() {
	ctx, recorder := s.newContext(http.MethodPut, `{"title":"new title"}`)
	ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

	s.controller.UpdateTask(ctx)

	s.Equal(http.StatusUnauthorized, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "UpdateTask", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestUpdateTask_StatusMapping() {
	updateReq := dto.TaskUpdateRequest{Title: "new title"}

	testCases := []struct {
		name       string
		mockErr    error
		wantStatus int
	}{
		{
			name:       "unauthorized owner",
			mockErr:    Usecases.ErrUnauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "task not found",
			mockErr:    errors.New("not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			mockUsecase := new(Mocks.MockTaskUsecase)
			controller := controllers.NewTaskController(mockUsecase)
			mockUsecase.On("UpdateTask", mock.Anything, s.taskID, s.ownerID, updateReq).Return(dto.TaskResponse{}, tc.mockErr)

			ctx, recorder := s.newContext(http.MethodPut, `{"title":"new title"}`)
			ctx.Set("userID", s.ownerID)
			ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

			controller.UpdateTask(ctx)

			s.Equal(tc.wantStatus, recorder.Code)
			mockUsecase.AssertExpectations(s.T())
		})
	}
}

func (s *TaskControllerSuite) TestDeleteTask_Success() {
	s.mockUsecase.On("DeleteTask", mock.Anything, s.taskID, s.ownerID).Return(nil)

	ctx, recorder := s.newContext(http.MethodDelete, "")
	ctx.Set("userID", s.ownerID)
	ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

	s.controller.DeleteTask(ctx)
	ctx.Writer.WriteHeaderNow()

	s.Equal(http.StatusNoContent, recorder.Code)
	s.Equal("", recorder.Body.String())

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TaskControllerSuite) TestDeleteTask_Unauthenticated() {
	ctx, recorder := s.newContext(http.MethodDelete, "")
	ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

	s.controller.DeleteTask(ctx)

	s.Equal(http.StatusUnauthorized, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "DeleteTask", mock.Anything, mock.Anything, mock.Anything)
}

func (s *TaskControllerSuite) TestDeleteTask_StatusMapping() {
	testCases := []struct {
		name       string
		mockErr    error
		wantStatus int
	}{
		{
			name:       "unauthorized owner",
			mockErr:    Usecases.ErrUnauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "task not found",
			mockErr:    errors.New("not found"),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			mockUsecase := new(Mocks.MockTaskUsecase)
			controller := controllers.NewTaskController(mockUsecase)
			mockUsecase.On("DeleteTask", mock.Anything, s.taskID, s.ownerID).Return(tc.mockErr)

			ctx, recorder := s.newContext(http.MethodDelete, "")
			ctx.Set("userID", s.ownerID)
			ctx.Params = gin.Params{{Key: "id", Value: s.taskID}}

			controller.DeleteTask(ctx)

			s.Equal(tc.wantStatus, recorder.Code)
			mockUsecase.AssertExpectations(s.T())
		})
	}
}
