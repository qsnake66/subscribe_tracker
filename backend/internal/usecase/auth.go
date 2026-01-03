package usecase

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"subscribe_tracker/backend/internal/domain"
)

type AuthUsecase struct {
	Users  UserRepository
	Tokens TokenManager
}

func NewAuthUsecase(users UserRepository, tokens TokenManager) AuthUsecase {
	return AuthUsecase{
		Users:  users,
		Tokens: tokens,
	}
}

type AuthResult struct {
	Token string
	User  domain.User
}

func (u AuthUsecase) Register(ctx context.Context, name, email, password string) (AuthResult, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	if name == "" || email == "" || len(password) < 8 {
		return AuthResult{}, ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := u.Users.Create(ctx, name, email, string(passwordHash))
	if err != nil {
		return AuthResult{}, err
	}

	token, err := u.Tokens.Sign(user.ID, user.Email)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{Token: token, User: user}, nil
}

func (u AuthUsecase) Login(ctx context.Context, email, password string) (AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return AuthResult{}, ErrInvalidInput
	}

	user, err := u.Users.FindByEmail(ctx, email)
	if err != nil {
		return AuthResult{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return AuthResult{}, ErrUnauthorized
	}

	token, err := u.Tokens.Sign(user.ID, user.Email)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{Token: token, User: user}, nil
}
