// package services

// import (
// 	"Learn_Jenkins/config"
// 	"Learn_Jenkins/domain/dto"
// 	"Learn_Jenkins/domain/model"
// 	"Learn_Jenkins/repositories"
// 	"context"
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/gorm"
// )

// func setupTestDB(t *testing.T) *gorm.DB {
// 	db, err := config.InitTestDatabase()
// 	if err != nil {
// 		t.Fatalf("failed to connect to test database: %v", err)
// 	}

// 	if err := db.AutoMigrate(&model.User{}); err != nil {
// 		t.Fatalf("failed to migrate: %v", err)
// 	}

// 	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
// 		t.Fatalf("failed to truncate users: %v", err)
// 	}

// 	t.Cleanup(func() {
// 		sqlDB, _ := db.DB()
// 		_ = sqlDB.Close()
// 	})

// 	return db
// }

// func TestMain(m *testing.M) {
// 	err := config.CreateTestDatabase()
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}

// 	code := m.Run()

// 	// dropErr := config.DropTestDatabase()
// 	// if dropErr != nil {
// 	// 	panic(dropErr)
// 	// }

// 	os.Exit(code)
// }

// func TestUserService_CreateUser(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := repositories.NewUserRepository(db)
// 	svc := NewUserService(repo)

// 	ctx := context.Background()
// 	req := &dto.UserRequest{Username: "Arthur"}

// 	user, err := svc.CreateUser(ctx, req)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, user)
// 	assert.Equal(t, "Arthur", user.Username)
// }

// func TestUserService_FindUserByID(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := repositories.NewUserRepository(db)
// 	svc := NewUserService(repo)

// 	ctx := context.Background()
// 	db.Create(&model.User{Username: "TestUser"})

// 	user, err := svc.FindUserByID(ctx, 1)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, user)
// 	assert.Equal(t, "TestUser", user.Username)
// }

// func TestUserService_FindAllUsers(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := repositories.NewUserRepository(db)
// 	svc := NewUserService(repo)

// 	ctx := context.Background()
// 	db.Create(&model.User{Username: "User1"})
// 	db.Create(&model.User{Username: "User2"})

// 	users, err := svc.FindAllUsers(ctx)

// 	assert.NoError(t, err)
// 	assert.Len(t, users, 2)
// 	assert.Equal(t, "User1", users[0].Username)
// 	assert.Equal(t, "User2", users[1].Username)
// }

package services

import (
	"context"
	"errors"
	"testing"

	"Learn_Jenkins/domain/dto"
	"Learn_Jenkins/domain/model"

	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	createResp  *model.User
	createErr   error
	findResp    *model.User
	findErr     error
	findAllResp []*model.User
	findAllErr  error
}

func (m *mockUserRepo) CreateUser(ctx context.Context, req *dto.UserRequest) (*model.User, error) {
	return m.createResp, m.createErr
}

func (m *mockUserRepo) FindUserByID(ctx context.Context, id uint) (*model.User, error) {
	return m.findResp, m.findErr
}

func (m *mockUserRepo) FindAllUsers(ctx context.Context) ([]*model.User, error) {
	return m.findAllResp, m.findAllErr
}

func TestUserService_CreateUser_WithMock_Success(t *testing.T) {
	ctx := context.Background()
	mock := &mockUserRepo{
		createResp: &model.User{ID: 1, Username: "Arthur"},
	}
	svc := NewUserService(mock)

	resp, err := svc.CreateUser(ctx, &dto.UserRequest{Username: "Arthur"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "Arthur", resp.Username)
}

func TestUserService_CreateUser_WithMock_RepoError(t *testing.T) {
	ctx := context.Background()
	mock := &mockUserRepo{
		createErr: errors.New("db error"),
	}
	svc := NewUserService(mock)

	resp, err := svc.CreateUser(ctx, &dto.UserRequest{Username: "Arthur"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUserService_FindUserByID_WithMock_Success(t *testing.T) {
	ctx := context.Background()
	mock := &mockUserRepo{
		findResp: &model.User{ID: 2, Username: "TestUser"},
	}
	svc := NewUserService(mock)

	resp, err := svc.FindUserByID(ctx, 2)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(2), resp.ID)
	assert.Equal(t, "TestUser", resp.Username)
}

func TestUserService_FindUserByID_WithMock_RepoError(t *testing.T) {
	ctx := context.Background()
	mock := &mockUserRepo{
		findErr: errors.New("not found"),
	}
	svc := NewUserService(mock)

	resp, err := svc.FindUserByID(ctx, 2)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUserService_FindAllUsers_WithMock_Success(t *testing.T) {
	ctx := context.Background()
	mock := &mockUserRepo{
		findAllResp: []*model.User{
			{ID: 1, Username: "User1"},
			{ID: 2, Username: "User2"},
		},
	}
	svc := NewUserService(mock)

	resp, err := svc.FindAllUsers(ctx)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "User1", resp[0].Username)
	assert.Equal(t, "User2", resp[1].Username)
}

func TestUserService_FindAllUsers_WithMock_RepoError(t *testing.T) {
	ctx := context.Background()
	mock := &mockUserRepo{
		findAllErr: errors.New("db failure"),
	}
	svc := NewUserService(mock)

	resp, err := svc.FindAllUsers(ctx)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
