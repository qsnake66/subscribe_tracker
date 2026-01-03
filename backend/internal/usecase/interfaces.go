package usecase

import (
	"context"

	"subscribe_tracker/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, name, email, passwordHash string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type SubscriptionRepository interface {
	ListByUserID(ctx context.Context, userID string) ([]domain.Subscription, error)
	Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	Delete(ctx context.Context, userID, id string) error
}

type TokenManager interface {
	Sign(userID, email string) (string, error)
	Parse(token string) (string, error)
}
