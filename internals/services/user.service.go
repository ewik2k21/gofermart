package services

import (
	"context"
	"gofermart/internals/interfaces"
	"gofermart/internals/repositories"
)

type IUserService interface {
}

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUserAccount(userRequest *interfaces.UserRequest, ctx context.Context) (string, error) {
	userId, err := s.userRepo.CreateUserAccount(userRequest, ctx)
	if err != nil {
		return "", err
	}

	return userId, nil
}
