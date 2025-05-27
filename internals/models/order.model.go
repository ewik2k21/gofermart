package models

import "github.com/gofrs/uuid"

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
	OrderNumber int         `db:"order_number"`
	Status      OrderStatus `db:"status"`
}
