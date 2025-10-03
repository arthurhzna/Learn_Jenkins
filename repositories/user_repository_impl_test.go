package repositories

import (
	"Learn_Jenkins/config"
	"Learn_Jenkins/domain/dto"
	"Learn_Jenkins/domain/model"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := config.InitTestDatabase()
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("failed to truncate users: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})

	return db
}

func TestMain(m *testing.M) {
	err := config.CreateTestDatabase()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	code := m.Run()

	// dropErr := config.DropTestDatabase()
	// if dropErr != nil {
	// 	panic(dropErr)
	// }

	os.Exit(code)
}

func TestUserRepository_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	ctx := context.Background()
	req := &dto.UserRequest{Username: "Arthur"}

	user, err := repo.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Arthur", user.Username)
}

func TestUserRepository_FindUserByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	ctx := context.Background()
	db.Create(&model.User{Username: "TestUser"})

	user, err := repo.FindUserByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "TestUser", user.Username)
}

func TestUserRepository_FindAllUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	ctx := context.Background()
	db.Create(&model.User{Username: "User1"})
	db.Create(&model.User{Username: "User2"})

	users, err := repo.FindAllUsers(ctx)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "User1", users[0].Username)
	assert.Equal(t, "User2", users[1].Username)
}
