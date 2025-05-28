package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type OrderStatus string

const (
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusNew        OrderStatus = "NEW"
)

type Order struct {
	ID          uuid.UUID   `db:"id"`
	UserID      uuid.UUID   `db:"user_id"`
	OrderNumber string      `db:"order_number"`
	Status      OrderStatus `db:"status"`
	Accrual     int         `db:"accrual"`
	UpdateAt    time.Time   `db:"update_at"`
}
