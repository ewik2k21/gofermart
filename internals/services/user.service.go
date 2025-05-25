package services

import (
	"context"
	"gofermart/internals/interfaces"
	"gofermart/internals/repositories"
)

type IUserService interface {
	CreateUserAccount(userRequest *interfaces.UserRequest, ctx context.Context) (string, error)
}

type UserService struct {
	userRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) IUserService {
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
