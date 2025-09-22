package service

import (
	"context"
	"fmt"
	"subscription_service/pkg/model"
	"subscription_service/pkg/repository"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreateSubscription(ctx context.Context, req *model.CreateSubscriptionRequest) error
	ListSubscriptions(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error)
}

type subService struct {
	repo   repository.Repository
	logger *logrus.Logger
}

func NewSubService(repo repository.Repository, logger *logrus.Logger) Service {
	return &subService{repo: repo, logger: logger}
}

func (s *subService) CreateSubscription(ctx context.Context, req *model.CreateSubscriptionRequest) error {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		s.logger.WithError(err).Error("Неверный формат даты начала подписки")
		return fmt.Errorf("неверный формат даты начала подписки, ожидается формат MM-YYYY")
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			s.logger.WithError(err).Error("Неверный формат даты окончания подписки")
			return fmt.Errorf("неверный формат даты окончания подписки, ожидается формат MM-YYYY")
		}
		endDate = &parsedEndDate
	}

	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		return err
	}

	return nil
}

func (s *subService) ListSubscriptions(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*model.Subscription, error) {
	return s.repo.List(ctx, userID, serviceName)
}
