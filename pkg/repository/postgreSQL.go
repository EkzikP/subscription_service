package repository

import (
	"context"
	"fmt"
	"subscription_service/config"

	"github.com/jackc/pgx/v5"
)

func NewPostgresDB(config *config.PQ) (pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.User, config.Pass, config.Host, config.Port, config.Name)
	db, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return pgx.Conn{}, err
	}
	return *db, nil
}
