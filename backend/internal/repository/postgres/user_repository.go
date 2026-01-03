package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"subscribe_tracker/backend/internal/domain"
	"subscribe_tracker/backend/internal/usecase"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return UserRepository{DB: db}
}

func (r UserRepository) Create(ctx context.Context, name, email, passwordHash string) (domain.User, error) {
	var user domain.User
	err := r.DB.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password_hash
	`, name, email, passwordHash).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return domain.User{}, usecase.ErrEmailExists
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.DB.QueryRow(ctx, `
		SELECT id, name, email, password_hash
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, usecase.ErrUnauthorized
		}
		return domain.User{}, err
	}
	return user, nil
}
