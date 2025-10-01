package controllers

import (
	"Learn_Jenkins/domain/dto"
	"Learn_Jenkins/services"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeUserService struct {
	createResp  *dto.UserResponse
	createErr   error
	findResp    *dto.UserResponse
	findErr     error
	findAllResp []*dto.UserResponse
	findAllErr  error
}

func (f *fakeUserService) CreateUser(ctx context.Context, req *dto.UserRequest) (*dto.UserResponse, error) {
	return f.createResp, f.createErr
}

func (f *fakeUserService) FindUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	return f.findResp, f.findErr
}

func (f *fakeUserService) FindAllUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	return f.findAllResp, f.findAllErr
}

func TestUserController_CreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fake := &fakeUserService{
		createResp: &dto.UserResponse{ID: 1, Username: "Arthur"},
	}
	ctrl := NewUserController(services.UserService(fake))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username":"Arthur"}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ctrl.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Arthur", resp.Username)
}

func TestUserController_CreateUser_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fake := &fakeUserService{}
	ctrl := NewUserController(services.UserService(fake))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// missing username -> validation should fail
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ctrl.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_CreateUser_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fake := &fakeUserService{
		createErr: errors.New("service failure"),
	}
	ctrl := NewUserController(services.UserService(fake))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username":"Arthur"}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ctrl.CreateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserController_FindUserByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fake := &fakeUserService{
		findResp: &dto.UserResponse{ID: 1, Username: "TestUser"},
	}
	ctrl := NewUserController(services.UserService(fake))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// set param id
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	ctrl.FindUserByID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "TestUser", resp.Username)
}

func TestUserController_FindUserByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fake := &fakeUserService{}
	ctrl := NewUserController(services.UserService(fake))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	ctrl.FindUserByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_FindAllUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fake := &fakeUserService{
		findAllResp: []*dto.UserResponse{
			{ID: 1, Username: "User1"},
			{ID: 2, Username: "User2"},
		},
	}
	ctrl := NewUserController(services.UserService(fake))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctrl.FindAllUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []*dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "User1", resp[0].Username)
	assert.Equal(t, "User2", resp[1].Username)
}
