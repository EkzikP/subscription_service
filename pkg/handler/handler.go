package handler

import (
	"net/http"
	"subscription_service/pkg/model"
	"subscription_service/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service service.Service
	logger  *logrus.Logger
}

func NewSubHandler(service service.Service, logger *logrus.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// CreateSubscription godoc
// @Summary Создать новую подписку
// @Description Добавляет новую запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.CreateSubscriptionRequest true "Subscription details"
// @Success 201 {object} nil
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(ctx *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Неверное тело запроса")
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Неверное тело запроса"})
		return
	}

	err := h.service.CreateSubscription(ctx.Request.Context(), &req)
	if err != nil && err.Error() == "запись уже существует" {
		h.logger.Warn("Ошибка добавления записи, запись уже существует")
		ctx.JSON(http.StatusConflict, model.ErrorResponse{Error: "Подписка уже существует для данного пользователя"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
	}

	ctx.Status(http.StatusCreated)
}

// ListSubscriptions godoc
// @Summary Список подписок
// @Description Получить список подписок с возможностью фильтрации по пользователю и сервису
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service Name"
// @Success 200 {array} model.Subscription
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(ctx *gin.Context) {
	var userID *uuid.UUID
	var serviceName *string

	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		parsedUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			h.logger.WithError(err).Warn("Неверный формат UserID")
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Неверный формат UserID"})
			return
		}
		userID = &parsedUUID
	}

	if serviceNameStr := ctx.Query("service_name"); serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	subscriptions, err := h.service.ListSubscriptions(ctx.Request.Context(), userID, serviceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subscriptions)
}
