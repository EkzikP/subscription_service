package handler

import (
	"net/http"
	"subscription_service/pkg/model"
	"subscription_service/pkg/service"

	"github.com/gin-gonic/gin"
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
