package main

import (
	"context"
	"fmt"
	"os"
	"subscription_service/pkg/config"
	"subscription_service/pkg/handler"
	"subscription_service/pkg/service"
	"time"

	_ "subscription_service/docs"
	"subscription_service/pkg/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Subscriptions Service API
// @version         1.0
// @description     REST‑сервис для учёта онлайн‑подписок пользователей.
// @host localhost:8080
// @BasePath        /

var logger = logrus.New()
var cfg = config.LoadConfig(logger)

func main() {

	// Инициализация логгера
	openLogFiles, err := InitLogrus()
	if err != nil {
		logger.Error("Ошибка при настройке параметров логгера. Вывод всех ошибок будет осуществлён в консоль")
	} else {
		// Закрытие всех открытых файлов в результате настройки логгера
		defer openLogFiles.Close()
	}

	// Создание нового подключения к БД
	pool, connStr, err := connectToDB()
	if err != nil {
		logger.Fatal("Ошибка подключения к БД: ", err)
	}
	defer pool.Close()

	// Миграция БД
	err = repository.MigrateDB(connStr)
	if err != nil {
		if err.Error() == "no change" {
			logger.Info("Миграция БД не требуется")
		} else {
			logger.Fatal("Ошибка миграции БД: ", err)
		}
	} else {
		logger.Info("Выполнена миграция БД")
	}

	// Инициализация зависимостей
	subRepo := repository.NewSubRepo(pool, logger)
	subService := service.NewSubService(subRepo, logger)
	subHandler := handler.NewSubHandler(subService, logger)

	r := setRouter(subHandler)

	logger.Info("Запуск сервера на порту ", cfg.HTTPPort)
	err = r.Run(":" + cfg.HTTPPort)
	if err != nil {
		logger.Fatal("Ошибка запуска сервера: ", err)
	}
}

func InitLogrus() (file *os.File, err error) {
	// установим уровень логирования
	switch cfg.LogLevel {
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// установим форматирование логов в джейсоне
	logger.SetFormatter(&logrus.JSONFormatter{})

	// установим вывод логов в файл
	file, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
		return file, nil
	} else {
		logger.Warn("Не удалось открыть файл логов, используется стандартный stderr")
		return nil, err
	}
}

func connectToDB() (*pgxpool.Pool, string, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	confPool, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, confPool)
	if err != nil {
		return nil, "", err
	}

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		return nil, "", err
	}

	logger.Info("Выполнено подключение к БД")
	return pool, connStr, nil
}

func setRouter(handler *handler.Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	rout := gin.New()
	rout.Use(gin.Recovery())
	rout.Use(loggingMiddleware())

	// Swagger
	rout.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Маршруты
	sub := rout.Group("/subscriptions")
	{
		sub.POST("", handler.CreateSubscription)
		sub.GET("", handler.ListSubscriptions)
		sub.GET("/:user_id/:service_name", handler.GetSubscription)
		sub.PUT("/:user_id/:service_name", handler.UpdateSubscription)
		sub.DELETE("/:user_id/:service_name", handler.DeleteSubscription)
		sub.GET("/total", handler.GetTotalSubscriptions)
	}
	return rout
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP request")
	}
}
