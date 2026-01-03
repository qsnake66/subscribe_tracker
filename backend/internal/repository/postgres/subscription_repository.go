package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"subscribe_tracker/backend/internal/domain"
	"subscribe_tracker/backend/internal/usecase"
)

type SubscriptionRepository struct {
	DB *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) SubscriptionRepository {
	return SubscriptionRepository{DB: db}
}

func (r SubscriptionRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Subscription, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, user_id, service_name, bank_name, card_last4, billing_cycle, charge_date
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY charge_date ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Subscription
	for rows.Next() {
		var item domain.Subscription
		var chargeDate time.Time
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ServiceName,
			&item.BankName,
			&item.CardLast4,
			&item.Billing,
			&chargeDate,
		); err != nil {
			return nil, err
		}
		item.ChargeDate = chargeDate
		results = append(results, item)
	}
	return results, nil
}

func (r SubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	var created domain.Subscription
	var chargeDate time.Time
	err := r.DB.QueryRow(ctx, `
		INSERT INTO subscriptions (user_id, service_name, bank_name, card_last4, billing_cycle, charge_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, service_name, bank_name, card_last4, billing_cycle, charge_date
	`, sub.UserID, sub.ServiceName, sub.BankName, sub.CardLast4, sub.Billing, sub.ChargeDate).Scan(
		&created.ID,
		&created.UserID,
		&created.ServiceName,
		&created.BankName,
		&created.CardLast4,
		&created.Billing,
		&chargeDate,
	)
	if err != nil {
		return domain.Subscription{}, err
	}
	created.ChargeDate = chargeDate
	return created, nil
}

func (r SubscriptionRepository) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	var updated domain.Subscription
	var chargeDate time.Time
	err := r.DB.QueryRow(ctx, `
		UPDATE subscriptions
		SET service_name = $1, bank_name = $2, card_last4 = $3, billing_cycle = $4, charge_date = $5, updated_at = NOW()
		WHERE id = $6 AND user_id = $7
		RETURNING id, user_id, service_name, bank_name, card_last4, billing_cycle, charge_date
	`, sub.ServiceName, sub.BankName, sub.CardLast4, sub.Billing, sub.ChargeDate, sub.ID, sub.UserID).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.ServiceName,
		&updated.BankName,
		&updated.CardLast4,
		&updated.Billing,
		&chargeDate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Subscription{}, usecase.ErrNotFound
		}
		return domain.Subscription{}, err
	}
	updated.ChargeDate = chargeDate
	return updated, nil
}

func (r SubscriptionRepository) Delete(ctx context.Context, userID, id string) error {
	cmd, err := r.DB.Exec(ctx, `
		DELETE FROM subscriptions WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}
	return nil
}
