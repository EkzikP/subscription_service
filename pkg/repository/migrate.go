package repository

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDB(connString string) error {
	m, err := migrate.New(
		"file://migrations",
		connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}
	return nil
}
