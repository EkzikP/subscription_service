package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"subscription_service/pkg/model/dto"
	"subscription_service/pkg/model/entities"
	"subscription_service/pkg/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

// CreateSubscription godoc
// @Summary      Создать подписку
// @Description  Добавляет новую запись о подписке
// @Accept       json
// @Produce      json
// @Param        subscription  body      dto.SubscriptionPayload  true  "data"
// @Success      201  {object}  nil
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Router       /subscriptions [post]
func (h *Handler) CreateSubscription(ctx *gin.Context) {

	var req dto.Subscription

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "неверный формат user_id"})
		return
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "неверный формат start_date"})
		return
	}

	var endDate time.Time
	if req.EndDate != "" {
		endDate, err = parseDate(req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "неверный формат end_date"})
			return
		}
	}

	_, err = repository.Select(h.pool, req.ServiceName, uid)
	if err == nil {
		ctx.JSON(http.StatusConflict, dto.ErrorResponse{Error: "подписка уже существует для данного пользователя"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	sub := entities.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      uid,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	err = repository.AddSubscription(h.pool, sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("01-2006", date)
}
