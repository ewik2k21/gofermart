package repositories

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gofermart/internals/interfaces"
	"gofermart/internals/utils"
	"golang.org/x/net/context"
)

type IUserRepository interface {
	CreateUserAccount(request *interfaces.UserRequest, ctx context.Context) (string, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUserAccount(request *interfaces.UserRequest, ctx context.Context) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("failed uuid generate : %v", err)
		return "", err
	}
	salt := id.String()[:9]
	passwordHash := utils.GeneratePasswordHash(request.Password, salt)

	sqlQuery := `INSERT INTO users(id,login, password_hash)VALUES ($1,$2,$3) RETURNING id`

	err = r.db.QueryRow(sqlQuery, id, request.Login, passwordHash).Scan(&id)
	if err != nil {
		logrus.Errorf("failed create user : %v", err)
	}

	return id.String(), err

}
