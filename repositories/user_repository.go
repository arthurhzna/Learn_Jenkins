package repositories

import (
	"Learn_Jenkins/domain/dto"
	"Learn_Jenkins/domain/model"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *dto.UserRequest) (*model.User, error)
	FindUserByID(ctx context.Context, id uint) (*model.User, error)
	FindAllUsers(ctx context.Context) ([]*model.User, error)
}
