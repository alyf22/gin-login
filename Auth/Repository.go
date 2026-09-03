package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repository *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := repository.db.QueryRow(ctx, `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *AuthRepository) CreateUser(ctx context.Context, user *User) error {
	return repository.db.QueryRow(ctx, `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`, user.Name, user.Email, user.Password).Scan(&user.ID)
}

var _ UserAuthRepository = (*AuthRepository)(nil)

func IsUserNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
