package repository

import (
	"context"
	"database/sql"
	"fmt"
	"subscription_service/pkg/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error)
	Create(ctx context.Context, sub *model.Subscription) error
	Get(ctx context.Context, userID *uuid.UUID, serviceName *string) (*model.Subscription, error)
	Update(ctx context.Context, updSub *model.UpdateSubscription) (*model.Subscription, error)
	Delete(ctx context.Context, userID *uuid.UUID, serviceName *string) error
}

type repo struct {
	pool   *pgxpool.Pool
	logger *logrus.Logger
}

func NewSubRepo(pool *pgxpool.Pool, logger *logrus.Logger) Repository {
	return &repo{pool: pool, logger: logger}
}

func (r *repo) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error) {
	query := `SELECT service_name, price, user_id, start_date, end_date 
              FROM subscriptions WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *userID)
		argCount++
	}

	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argCount)
		args = append(args, *serviceName)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Ошибка при выполнении запроса списка подписок")
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(
			&sub.ServiceName, &sub.Price, &sub.UserID,
			&sub.StartDate, &sub.EndDate,
		)
		if err != nil {
			r.logger.WithError(err).Error("Ошибка при сканировании строки результата запроса")
			return nil, err
		}
		subscriptions = append(subscriptions, &sub)
	}

	r.logger.WithField("count", len(subscriptions)).Info("Получен список подписок")
	return subscriptions, nil
}

func (r *repo) Create(ctx context.Context, sub *model.Subscription) error {

	exist, err := r.List(ctx, &sub.UserID, &sub.ServiceName)
	if err != nil {
		return err
	}

	if len(exist) > 0 {
		return fmt.Errorf("запись уже существует")
	}

	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
              VALUES ($1, $2, $3, $4, $5) RETURNING user_id, service_name`

	err = r.pool.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).
		Scan(&sub.UserID, &sub.ServiceName)

	if err != nil {
		r.logger.WithError(err).Error("Ошибка при создании записи в таблице subscriptions")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"serviceName": sub.ServiceName,
		"userId":      sub.UserID,
	}).Info("Запись о подписке успешно создана")

	return nil
}

func (r *repo) Get(ctx context.Context, userID *uuid.UUID, serviceName *string) (*model.Subscription, error) {
	query := `SELECT service_name, price, user_id, start_date, end_date 
              FROM subscriptions WHERE user_id = $1 AND service_name = $2`

	var sub model.Subscription
	err := r.pool.QueryRow(ctx, query, userID, serviceName).
		Scan(&sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *repo) Update(ctx context.Context, updSub *model.UpdateSubscription) (*model.Subscription, error) {
	query := `UPDATE subscriptions SET 
                price = COALESCE($1, price),
                start_date = COALESCE($2, start_date),
                end_date = $3
              WHERE user_id = $4 AND service_name = $5 RETURNING *`

	var sub model.Subscription
	err := r.pool.QueryRow(ctx, query, updSub.Price, updSub.StartDate, updSub.EndDate, updSub.UserID, updSub.ServiceName).
		Scan(&sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *repo) Delete(ctx context.Context, userID *uuid.UUID, serviceName *string) error {
	query := "DELETE FROM subscriptions WHERE user_id = $1 AND service_name = $2"

	result, err := r.pool.Exec(ctx, query, userID, serviceName)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		err := sql.ErrNoRows
		return err
	}

	return nil
}
