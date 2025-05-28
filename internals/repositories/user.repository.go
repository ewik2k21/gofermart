package repositories

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gofermart/internals/interfaces"
	"gofermart/internals/models"
	"golang.org/x/net/context"
)

type IUserRepository interface {
	CreateUserAccount(request *interfaces.UserRequest, id uuid.UUID, ctx context.Context) error
	GetUserByLogin(login string) (*interfaces.UserLoginData, error)
	AddOrder(userId, orderNumber string) (string, error)
	GetAllOrders(userId string) (*[]models.Order, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUserAccount(request *interfaces.UserRequest, id uuid.UUID, ctx context.Context) error {

	sqlQuery := `INSERT INTO users(id,login, password_hash)VALUES ($1,$2,$3) RETURNING id`

	_, err := r.db.Exec(sqlQuery, id, request.Login, request.Password)
	if err != nil {
		logrus.Errorf("failed create user : %v", err)
		return err
	}

	return nil

}

func (r *UserRepository) GetUserByLogin(login string) (*interfaces.UserLoginData, error) {
	sqlQuery := `SElECT id, password_hash FROM users WHERE login = $1 `
	userLoginData := &interfaces.UserLoginData{}

	err := r.db.QueryRow(sqlQuery, login).Scan(&userLoginData.UserId, &userLoginData.PasswordHash)
	if err != nil {
		logrus.Errorf("failed get user by login : %v", err)
		return nil, err
	}

	return userLoginData, nil
}

func (r *UserRepository) AddOrder(userId, orderNumber string) (string, error) {
	sqlQuery := `
	with ExistingOrder as (
	    select user_id 
	    from orders
	    where order_number = $1
	), 
	InsertOrder as (
	    insert into orders (order_number, user_id, status)
	    select $1,$2,$3
	    where not exists (select 1 from ExistingOrder)
	)
	select user_id from ExistingOrder`
	var id uuid.UUID
	err := r.db.QueryRow(sqlQuery, orderNumber, userId, models.OrderStatusNew).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return id.String(), nil
}

func (r *UserRepository) GetAllOrders(userId string) (*[]models.Order, error) {
	orders := &[]models.Order{}

	sqlQuery := `
	select order_number,status,update_at, accrual 
	from orders
	where user_id = $1 
	order by update_at`

	rows, err := r.db.Query(sqlQuery, userId)
	if err != nil {
		logrus.Errorf("failed get all orders : %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.OrderNumber, &order.Status, &order.UpdateAt, &order.Accrual)
		if err != nil {
			logrus.Errorf("failed scan row: %v", err)
			return nil, err
		}
		*orders = append(*orders, order)
	}
	if err := rows.Err(); err != nil {
		logrus.Errorf("failed iteration on rows: %v", err)
		return nil, err
	}

	return orders, nil
}
