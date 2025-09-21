package repository

import (
	"context"
	"subscription_service/pkg/model/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddSubscription(pool *pgxpool.Pool, subscription entities.Subscription) error {
	_, err := pool.Exec(context.Background(),
		"INSERT INTO service.subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)",
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate)
	if err != nil {
		return err
	}
	return nil
}
