package repository

import (
	"context"
	"fmt"
	"subscription_service/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(config *config.Config) (*pgxpool.Pool, error) {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Db.User,
		config.Db.Pass,
		config.Db.Host,
		config.Db.Port,
		config.Db.Name,
		config.Db.Ssl)

	config.ConString = connString

	pqConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pqConfig)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func MigrateDB(connString string) error {
	m, err := migrate.New(
		"file://pkg/repository/migrations",
		connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}
	return nil
}
