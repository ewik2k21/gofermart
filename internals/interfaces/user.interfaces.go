package interfaces

import (
	"github.com/gofrs/uuid"
	"gofermart/internals/models"
)

type UserRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserData struct {
	ID      uuid.UUID      `json:"id"`
	Login   string         `json:"login"`
	Balance models.Balance `json:"balance"`
	Orders  []models.Order `json:"orders"`
}

type UserLoginData struct {
	UserId       uuid.UUID
	PasswordHash string
}

type OrderRequest struct {
	OrderNumber string `json:"order_number" validate:"required"`
}
