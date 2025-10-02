package services

import (
	"Learn_Jenkins/domain/dto"
	"Learn_Jenkins/repositories"
	"context"
)

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userServiceImpl{userRepository: userRepository}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req *dto.UserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepository.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *userServiceImpl) FindUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *userServiceImpl) FindAllUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := s.userRepository.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	var responses []*dto.UserResponse
	for _, user := range users {
		responses = append(responses, &dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
		})
	}
	return responses, nil
}
