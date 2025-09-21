package main

import (
	"log"
	"subscription_service/pkg/config"
	"subscription_service/pkg/router"

	//	_ "subscription_service_my/docs"
	"subscription_service/pkg/repository"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = logrus.New()

// @title           Subscriptions Service API
// @version         1.0
// @description     REST‑сервис для учёта онлайн‑подписок пользователей.
// @host localhost:8080
// @BasePath        /

func init() {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Print("Не найден .env файл")
	}
}

func main() {

	// Инициализация переменных внешней среды
	cfg := config.New()

	// Инициализация логгера
	openLogFiles, err := config.InitLogrus(cfg.LogLevel, logger)
	if err != nil {
		logger.Error("Ошибка при настройке параметров логгера. Вывод всех ошибок будет осуществлён в консоль")
	} else {
		// Закрытие всех открытых файлов в результате настройки логгера
		defer openLogFiles.Close()
	}

	// Создание нового подключения к БД
	pool, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Выполнено подключение к БД")

	defer pool.Close()

	// Миграция БД
	err = repository.MigrateDB(cfg.ConString)
	if err != nil {
		if err.Error() == "no change" {
			logger.Info("Миграция БД не требуется")
		} else {
			logger.Fatal(err)
		}
	} else {
		logger.Info("Выполнена миграция БД")
	}

	r := router.SetRouter(pool, logger)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = r.Run(":" + cfg.HttpPort)
	if err != nil {
		logger.Fatal(err)
	}
}
