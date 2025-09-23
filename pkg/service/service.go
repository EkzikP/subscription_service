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
	GetSubscription(ctx context.Context, userID *uuid.UUID, serviceName *string) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, userID *uuid.UUID, serviceName *string, req *model.UpdateSubscriptionRequest) (*model.Subscription, error)
	DeleteSubscription(ctx context.Context, userID *uuid.UUID, serviceName *string) error
}

type subService struct {
	repo   repository.Repository
	logger *logrus.Logger
}

func NewSubService(repo repository.Repository, logger *logrus.Logger) Service {
	return &subService{repo: repo, logger: logger}
}

func (s *subService) CreateSubscription(ctx context.Context, req *model.CreateSubscriptionRequest) error {

	startDate, endDate, err := ParseDate(&req.StartDate, req.EndDate, s.logger)
	if err != nil {
		return err
	}

	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   *startDate,
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

func (s *subService) GetSubscription(ctx context.Context, userID *uuid.UUID, serviceName *string) (*model.Subscription, error) {
	return s.repo.Get(ctx, userID, serviceName)
}

func (s *subService) UpdateSubscription(ctx context.Context, userID *uuid.UUID, serviceName *string, req *model.UpdateSubscriptionRequest) (*model.Subscription, error) {

	startDate, endDate, err := ParseDate(req.StartDate, req.EndDate, s.logger)
	if err != nil {
		return nil, err
	}

	updSub := &model.UpdateSubscription{
		ServiceName: *serviceName,
		UserID:      *userID,
		Price:       req.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return s.repo.Update(ctx, updSub)
}

func (s *subService) DeleteSubscription(ctx context.Context, userID *uuid.UUID, serviceName *string) error {
	return s.repo.Delete(ctx, userID, serviceName)
}

func ParseDate(startDayStr *string, endDayStr *string, logger *logrus.Logger) (*time.Time, *time.Time, error) {
	var startDate *time.Time
	var endDate *time.Time

	if startDayStr != nil {
		parsedStartDate, err := time.Parse("01-2006", *startDayStr)
		if err != nil {
			logger.WithError(err).Error("Неверный формат даты начала подписки")
			return startDate, endDate, fmt.Errorf("неверный формат даты начала подписки, ожидается формат MM-YYYY")
		}
		startDate = &parsedStartDate
	}

	if endDayStr != nil {
		parsedEndDate, err := time.Parse("01-2006", *endDayStr)
		if err != nil {
			logger.WithError(err).Error("Неверный формат даты окончания подписки")
			return startDate, endDate, fmt.Errorf("неверный формат даты окончания подписки, ожидается формат MM-YYYY")
		}
		endDate = &parsedEndDate
	}

	return startDate, endDate, nil
}
