package repository

import (
	"context"
	"time"

	"github.com/Wookkie/subscription-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) CreateSubscription(sub domain.Subscription) (*domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, `
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
        VALUES ($1,$2,$3,$4,$5,$6)
        RETURNING id`,
		sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(&sub.ID)

	return &sub, err
}

func (r *SubscriptionRepo) GetAllSubscriptions() ([]domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sub []domain.Subscription

	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate); err != nil {
			return nil, err
		}
		sub = append(sub, s)
	}

	return sub, nil
}

func (r *SubscriptionRepo) GetSubscriptionByID(id string) (*domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sub domain.Subscription

	err := r.db.QueryRow(ctx, `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions WHERE id=$1`, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)

	return &sub, err
}

func (r *SubscriptionRepo) UpdateSubscription(id string, sub domain.Subscription) (*domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, `
        UPDATE subscriptions
        SET price=$1, end_date=$2
        WHERE id=$3
        RETURNING id, service_name, price, user_id, start_date, end_date`,
		sub.Price, sub.EndDate, id,
	).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)

	return &sub, err
}

func (r *SubscriptionRepo) DeleteSubscription(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `DELETE FROM subscriptions WHERE id=$1`, id)

	return err
}

func (r *SubscriptionRepo) GetSubscriptionsForPeriod(from, to time.Time, userID, serviceName string) ([]domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions
        WHERE start_date >= $1 AND start_date <= $2
    `

	args := []any{from, to}

	if userID != "" {
		query += " AND user_id = $3"
		args = append(args, userID)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sub []domain.Subscription

	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate); err != nil {
			return nil, err
		}
		sub = append(sub, s)
	}

	return sub, nil
}
