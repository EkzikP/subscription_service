package repository

import (
	"context"
	"subscription_service/pkg/model/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Select(pool *pgxpool.Pool, service_name string, user_id uuid.UUID) (entities.Subscription, error) {
	row := pool.QueryRow(context.Background(),
		"SELECT * FROM service.subscriptions WHERE service_name = $1 AND user_id = $2",
		service_name,
		user_id)
	var subscription entities.Subscription
	if err := row.Scan(&subscription); err != nil {
		return subscription, err
	}

	return subscription, nil
}
