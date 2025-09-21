package router

import (
	"subscription_service/pkg/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
)

func SetRouter(pool *pgxpool.Pool, logger *logrus.Logger) *gin.Engine {

	r := gin.New()
	r.Use(ginlogrus.Logger(logger), gin.Recovery())

	h := handlers.NewHandler(pool)

	r.POST("/subscriptions", h.CreateSubscription)
	//r.GET("/subscriptions", h.ListSubscriptions)
	//
	//r.GET("/subscriptions/:user_id/:service_name", h.GetSubscription)
	//r.PUT("/subscriptions/:user_id/:service_name", h.UpdateSubscription)
	//r.DELETE("/subscriptions/:user_id/:service_name", h.DeleteSubscription)
	//
	//r.GET("/subscriptions/summary", h.SumSubscriptions)

	return r

}
