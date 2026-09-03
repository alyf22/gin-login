package Auth

import (
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
}

type AuthRepository struct {
	db *pgx.Conn
}
