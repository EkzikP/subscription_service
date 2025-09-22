package handler

import (
	"database/sql"
	"errors"
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
	if err != nil {
		if err.Error() == "запись уже существует" {
			h.logger.Warn("Ошибка добавления записи, запись уже существует")
			ctx.JSON(http.StatusConflict, model.ErrorResponse{Error: "Подписка уже существует для данного пользователя"})
			return
		}
		h.logger.WithError(err).Error("Ошибка добавления записи")
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
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

	userID, serviceName, err := getUserIDAndServiceName(ctx)
	if err != nil {
		h.logger.WithError(err).Warn("Неверный формат UserID")
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Неверный формат UserID"})
		return
	}

	subscriptions, err := h.service.ListSubscriptions(ctx.Request.Context(), userID, serviceName)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка получения списка подписок")
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subscriptions)
}

// GetSubscription godoc
// @Summary Получить подписку
// @Description Получить подписку по ID пользователя и имени сервиса
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param service_name path string true "Имя сервиса"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(ctx *gin.Context) {

	userID, serviceName, err := getUserIDAndServiceName(ctx)
	if err != nil {
		h.logger.WithError(err).Warn("Неверный формат UserID")
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Неверный формат UserID"})
		return
	}

	subscription, err := h.service.GetSubscription(ctx.Request.Context(), userID, serviceName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.WithError(err).Warn("Подписка не найдена")
			ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Подписка не найдена"})
		}
		h.logger.WithError(err).Error("Ошибка получения подписки")
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
	}

	ctx.JSON(http.StatusOK, subscription)
}

func getUserIDAndServiceName(ctx *gin.Context) (*uuid.UUID, *string, error) {
	var userID *uuid.UUID
	var serviceName *string

	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		parsedUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, nil, err
		}
		userID = &parsedUUID
	}

	if serviceNameStr := ctx.Query("service_name"); serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	return userID, serviceName, nil
}
