package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required,min=1"`
	UserID      uuid.UUID  `json:"user_id" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	Price     *int    `json:"price,omitempty"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
}

type UpdateSubscription struct {
	ServiceName string     `json:"service_name" binding:"required"`
	UserID      uuid.UUID  `json:"user_id" binding:"required"`
	Price       *int       `json:"price,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type TotalRequest struct {
	UserID      *uuid.UUID `form:"user_id"`
	ServiceName *string    `form:"service_name"`
	StartPeriod string     `form:"start_period" binding:"required"`
	EndPeriod   string     `form:"end_period" binding:"required"`
}

type TotalSubscription struct {
	TotalCost int `json:"total_cost"`
	Count     int `json:"count"`
}
