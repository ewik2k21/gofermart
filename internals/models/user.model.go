package models

import "github.com/gofrs/uuid"

type User struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	Balance      Balance   `db:"balance"`
	Orders       []Order   `db:"orders"`
}

type Balance struct {
	Current  float64 `db:"current"`
	Withdraw float64 `db:"withdraw" default:"0.00"`
}
