package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"task_manager/Delivery/controllers"
	"task_manager/Domain/dto"
	"task_manager/Mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControllerSuite struct {
	suite.Suite
	mockUsecase *Mocks.MockUserUsecase
	controller  *controllers.UserController
}

func (s *UserControllerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockUsecase = new(Mocks.MockUserUsecase)
	s.controller = controllers.NewUserController(s.mockUsecase)
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerSuite))
}

func (s *UserControllerSuite) newContext(method, body string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(method, "/", strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx, recorder
}

func (s *UserControllerSuite) TestRegister_Success() {
	req := dto.RegisterRequest{Name: "nafiad", Email: "nafiad@example.com", Password: "password123"}
	resp := dto.UserResponse{ID: "64c3f5e6c9e13a2f3b30a001", Name: "nafiad", Email: "nafiad@example.com", Role: "user"}

	s.mockUsecase.On("Register", mock.Anything, req).Return(resp, nil)

	ctx, recorder := s.newContext(http.MethodPost, `{"name":"nafiad","email":"nafiad@example.com","password":"password123"}`)

	s.controller.Register(ctx)

	s.Equal(http.StatusCreated, recorder.Code)
	s.Contains(recorder.Body.String(), "nafiad")

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *UserControllerSuite) TestRegister_ValidationError() {
	ctx, recorder := s.newContext(http.MethodPost, `{"name":"nafiad","email":"not-an-email","password":"123"}`)

	s.controller.Register(ctx)

	s.Equal(http.StatusBadRequest, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "Register", mock.Anything, mock.Anything)
}

func (s *UserControllerSuite) TestRegister_ServiceFailure() {
	req := dto.RegisterRequest{Name: "nafiad", Email: "nafiad@example.com", Password: "password123"}

	s.mockUsecase.On("Register", mock.Anything, req).Return(dto.UserResponse{}, errors.New("Email already taken."))

	ctx, recorder := s.newContext(http.MethodPost, `{"name":"nafiad","email":"nafiad@example.com","password":"password123"}`)

	s.controller.Register(ctx)

	s.Equal(http.StatusBadRequest, recorder.Code)
	s.Contains(recorder.Body.String(), "Email already taken.")

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *UserControllerSuite) TestLogin_Success() {
	req := dto.LoginRequest{Email: "nafiad@example.com", Password: "password123"}

	s.mockUsecase.On("Login", mock.Anything, req).Return("signed-token", nil)

	ctx, recorder := s.newContext(http.MethodPost, `{"email":"nafiad@example.com","password":"password123"}`)

	s.controller.Login(ctx)

	s.Equal(http.StatusOK, recorder.Code)
	s.Contains(recorder.Body.String(), "signed-token")

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *UserControllerSuite) TestLogin_ValidationError() {
	ctx, recorder := s.newContext(http.MethodPost, `{"email":"nafiad@example.com"}`)

	s.controller.Login(ctx)

	s.Equal(http.StatusBadRequest, recorder.Code)

	s.mockUsecase.AssertNotCalled(s.T(), "Login", mock.Anything, mock.Anything)
}

func (s *UserControllerSuite) TestLogin_InvalidCredentials() {
	req := dto.LoginRequest{Email: "nafiad@example.com", Password: "wrong"}

	s.mockUsecase.On("Login", mock.Anything, req).Return("", errors.New("Invalid Credentials"))

	ctx, recorder := s.newContext(http.MethodPost, `{"email":"nafiad@example.com","password":"wrong"}`)

	s.controller.Login(ctx)

	s.Equal(http.StatusUnauthorized, recorder.Code)
	s.Contains(recorder.Body.String(), "Invalid Credentials.")

	s.mockUsecase.AssertExpectations(s.T())
}
