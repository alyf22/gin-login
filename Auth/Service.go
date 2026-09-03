package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	repository UserAuthRepository
	secret     []byte
	tokenTTL   time.Duration
}

func NewAuthService(repository UserAuthRepository, secret string) *AuthService {
	return &AuthService{
		repository: repository,
		secret:     []byte(secret),
		tokenTTL:   15 * time.Minute,
	}
}

func (service *AuthService) Login(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	if email == "" || request.Password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := service.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"email": user.Email,
		"iat": now.Unix(),
		"exp": now.Add(service.tokenTTL).Unix(),
	})

	accessToken, err := token.SignedString(service.secret)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{AccessToken: accessToken}, nil
}
