package services

import (
	"Learn_Jenkins/domain/dto"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.UserRequest) (*dto.UserResponse, error)
	FindUserByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	FindAllUsers(ctx context.Context) ([]*dto.UserResponse, error)
}
