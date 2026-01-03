package usecase

import (
	"context"
	"strings"
	"time"

	"subscribe_tracker/backend/internal/domain"
)

type SubscriptionUsecase struct {
	Subscriptions SubscriptionRepository
}

func NewSubscriptionUsecase(subscriptions SubscriptionRepository) SubscriptionUsecase {
	return SubscriptionUsecase{Subscriptions: subscriptions}
}

type SubscriptionInput struct {
	ServiceName string
	BankName    string
	CardLast4   string
	Billing     string
	ChargeDate  string
}

func (u SubscriptionUsecase) List(ctx context.Context, userID string) ([]domain.Subscription, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrUnauthorized
	}
	return u.Subscriptions.ListByUserID(ctx, userID)
}

func (u SubscriptionUsecase) Create(ctx context.Context, userID string, input SubscriptionInput) (domain.Subscription, error) {
	sub, err := u.toDomain(userID, input)
	if err != nil {
		return domain.Subscription{}, err
	}
	return u.Subscriptions.Create(ctx, sub)
}

func (u SubscriptionUsecase) Update(ctx context.Context, userID, id string, input SubscriptionInput) (domain.Subscription, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Subscription{}, ErrInvalidInput
	}
	sub, err := u.toDomain(userID, input)
	if err != nil {
		return domain.Subscription{}, err
	}
	sub.ID = id
	return u.Subscriptions.Update(ctx, sub)
}

func (u SubscriptionUsecase) Delete(ctx context.Context, userID, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidInput
	}
	return u.Subscriptions.Delete(ctx, userID, id)
}

func (u SubscriptionUsecase) toDomain(userID string, input SubscriptionInput) (domain.Subscription, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.Subscription{}, ErrUnauthorized
	}

	serviceName := strings.TrimSpace(input.ServiceName)
	bankName := strings.TrimSpace(input.BankName)
	cardLast4 := strings.TrimSpace(input.CardLast4)
	billing := strings.ToLower(strings.TrimSpace(input.Billing))
	if serviceName == "" || bankName == "" || len(cardLast4) != 4 {
		return domain.Subscription{}, ErrInvalidInput
	}

	if billing != "monthly" && billing != "yearly" {
		return domain.Subscription{}, ErrInvalidInput
	}

	chargeDate, err := time.Parse("2006-01-02", strings.TrimSpace(input.ChargeDate))
	if err != nil {
		return domain.Subscription{}, ErrInvalidInput
	}

	return domain.Subscription{
		UserID:      userID,
		ServiceName: serviceName,
		BankName:    bankName,
		CardLast4:   cardLast4,
		Billing:     billing,
		ChargeDate:  chargeDate,
	}, nil
}
