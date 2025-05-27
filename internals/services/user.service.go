package services

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gofermart/internals/interfaces"
	"gofermart/internals/repositories"
	"gofermart/internals/utils"
)

type IUserService interface {
	CreateUserAccount(userRequest *interfaces.UserRequest, ctx context.Context) (string, error)
	CheckCredentials(userRequest *interfaces.UserRequest) (string, bool, error)
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
	id, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("failed uuid generate : %v", err)
		return "", err
	}
	salt := id.String()[:9]
	passwordHash := utils.GeneratePasswordHash(userRequest.Password, salt)

	userRequest.Password = passwordHash
	err = s.userRepo.CreateUserAccount(userRequest, id, ctx)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *UserService) CheckCredentials(userRequest *interfaces.UserRequest) (string, bool, error) {

	userLoginData, err := s.userRepo.GetUserByLogin(userRequest.Login)
	if err != nil {
		return "", false, err
	}
	salt := userLoginData.UserId.String()[:9]
	ok := utils.DoPasswordMatch(userLoginData.PasswordHash, userRequest.Password, salt)

	return userLoginData.UserId.String(), ok, nil
}
