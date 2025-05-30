package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gofermart/internals/interfaces"
	"gofermart/internals/repositories"
	"gofermart/internals/utils"
)

type IUserService interface {
	CreateUserAccount(userRequest *interfaces.UserRequest, ctx context.Context) (string, error)
	CheckCredentials(userRequest *interfaces.UserRequest) (string, bool, error)
	AddOrder(userId, orderNumber string) (int, string, error)
	GetAllOrders(userId string) (*[]interfaces.OrderResponse, error)
	GetBalance(userId string) (*interfaces.BalanceResponse, error)
	FillDb(count int, ctx context.Context) error
	PostWithdraw(userId string, withdrawRequest interfaces.WithdrawRequest) (statusCode int, msg string, err error)
	GetWithdrawsById(userId string) (*[]interfaces.WithdrawsResponse, error)
}

type UserService struct {
	userRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) IUserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetWithdrawsById(userId string) (*[]interfaces.WithdrawsResponse, error) {
	withdraws, err := s.userRepo.GetWithdrawById(userId)
	if err != nil {
		return nil, err
	}
	if withdraws == nil {
		return nil, nil
	}

	var withdrawsResp []interfaces.WithdrawsResponse

	for _, withdraw := range *withdraws {
		withdrawsResp = append(withdrawsResp, interfaces.WithdrawsResponse{
			Order:       withdraw.Order,
			Sum:         withdraw.Sum,
			ProcessedAt: withdraw.ProcessedAt,
		})
	}
	return &withdrawsResp, nil
}

func (s *UserService) PostWithdraw(userId string, withdrawRequest interfaces.WithdrawRequest) (statusCode int, msg string, err error) {
	order, err := s.userRepo.GetOrderByNumberAndUserId(userId, withdrawRequest.Order)
	if err != nil && err != sql.ErrNoRows {
		return 500, "Failed: ...internal server error ", err
	}
	if order == "" || err == sql.ErrNoRows {
		return 422, "Incorrect order number", nil
	}
	balance, err := s.GetBalance(userId)
	if err != nil {
		return 500, "Failed: ...internal server error ", err
	}
	if float64(withdrawRequest.Sum) > balance.Current {
		return 402, "There are insufficient funds in the account", nil
	}
	err = s.userRepo.AddWithdraw(userId, withdrawRequest)
	if err != nil {
		return 500, "Failed: ...internal server error ", err
	}

	return 200, "Successfully add withdraws", nil

}

func (s *UserService) FillDb(count int, ctx context.Context) error {

	for i := 0; i < count; i++ {
		login := fmt.Sprintf("ewik2k%d", i)
		id, _ := uuid.NewV4()
		salt := id.String()[:9]
		passwordHash := utils.GeneratePasswordHash(login, salt)
		req := &interfaces.UserRequest{Login: login, Password: passwordHash}
		err := s.userRepo.CreateUserAccount(req, id, ctx)
		if err != nil {
			return err
		}
		orderNumber := fmt.Sprintf("123%d", i)
		_, err = s.userRepo.AddOrder(id.String(), orderNumber)
		if err != nil {
			return err
		}
		current := float64(i) * 2000
		withdraw := i * 1000
		err = s.userRepo.FillBalance(id, current, withdraw)
		if err != nil {
			return err
		}

	}

	return nil
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

func (s *UserService) AddOrder(userId, orderNumber string) (int, string, error) {
	userIdFromOrder, err := s.userRepo.AddOrder(userId, orderNumber)
	if userIdFromOrder == "" && err == nil {
		return 202, "New order number accepted for processing", nil
	}
	if err != nil {
		return 500, "Failed add order", err
	}
	if userIdFromOrder == userId {
		return 200, "Order number has already been uploaded by this user", nil
	}
	return 409, "Order number has already been uploaded by another user", nil

}

func (s *UserService) GetAllOrders(userId string) (*[]interfaces.OrderResponse, error) {
	orders, err := s.userRepo.GetAllOrders(userId)
	if err != nil {
		return nil, err
	}

	var ordersResp []interfaces.OrderResponse

	for _, order := range *orders {
		ordersResp = append(ordersResp, interfaces.OrderResponse{
			OrderNumber: order.OrderNumber,
			Status:      order.Status,
			Accrual:     order.Accrual,
			UpdatedAt:   order.UpdateAt,
		})
	}
	return &ordersResp, nil
}

func (s *UserService) GetBalance(userId string) (*interfaces.BalanceResponse, error) {
	balance, err := s.userRepo.GetBalance(userId)
	if err != nil {
		return nil, err
	}

	return &interfaces.BalanceResponse{Current: balance.Current, Withdrawn: balance.Withdraw}, nil
}
