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
	GetBalance(userId string) (*models.Balance, error)
	FillBalance(userId uuid.UUID, current float64, withdraw int) error
	GetOrderByNumberAndUserId(userId string, order string) (string, error)
	AddWithdraw(userId string, withdrawRequest interfaces.WithdrawRequest) error
	GetWithdrawById(userId string) (*[]models.Withdraw, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AddWithdraw(userId string, withdrawRequest interfaces.WithdrawRequest) error {
	updateOrderQuery := `UPDATE orders SET status = $1 WHERE order_number = $2`
	updateBalanceQuery := `UPDATE balances SET "current" = "current" - $1 WHERE user_id = $2 `
	insertQuery := `INSERT INTO withdraws(id,user_id, "order", sum) VALUES ($1,$2,$3,$4)`
	id, _ := uuid.NewV4()
	_, err := r.db.Exec(insertQuery, id, userId, withdrawRequest.Order, withdrawRequest.Sum)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(updateOrderQuery, models.OrderStatusProcessing, withdrawRequest.Order)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(updateBalanceQuery, float64(withdrawRequest.Sum), userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetOrderByNumberAndUserId(userId string, order string) (string, error) {
	sqlQuery := `SELECT order_number FROM orders WHERE user_id = $1 and order_number = $2`

	var orderNumber string

	err := r.db.QueryRow(sqlQuery, userId, order).Scan(&orderNumber)
	if err != nil && err == sql.ErrNoRows {
		return "", err
	}
	return orderNumber, nil

}

func (r *UserRepository) FillBalance(userId uuid.UUID, current float64, withdraw int) error {
	sqlQuery := `INSERT INTO balances (user_id, current, withdraw)
	VALUES ($1, $2, $3) `

	_, err := r.db.Exec(sqlQuery, userId, current, withdraw)
	if err != nil {
		return err
	}

	return nil

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

func (r *UserRepository) GetWithdrawById(userId string) (*[]models.Withdraw, error) {
	withdraws := &[]models.Withdraw{}

	sqlQuery := `
	select  "order", sum, processed_at 
	from withdraws
	where user_id = $1
	order by processed_at`

	rows, err := r.db.Query(sqlQuery, userId)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var withdraw models.Withdraw
		err = rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return nil, err
		}
		*withdraws = append(*withdraws, withdraw)
	}
	if err := rows.Err(); err != nil {
		logrus.Errorf("failed iteration on rows: %v", err)
		return nil, err
	}

	return withdraws, nil

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

func (r *UserRepository) GetBalance(userId string) (*models.Balance, error) {
	balance := &models.Balance{}

	sqlQuery := `
	select b.current, b.withdraw
	from balances as b 
	where user_id = $1`

	err := r.db.QueryRow(sqlQuery, userId).Scan(&balance.Current, &balance.Withdraw)
	if err != nil {
		logrus.Errorf("failed get balance : %v", err)
		return nil, err
	}

	return balance, nil
}
