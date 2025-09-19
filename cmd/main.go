package main

import (
	"log"
	"os"
	"subscription_service/config"
	_ "subscription_service/docs"
	"subscription_service/pkg/repository"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
	openLogFiles, err := config.InitLogrus(cfg.LogLevel)
	if err != nil {
		logrus.Error("Ошибка при настройке параметров логгера. Вывод всех ошибок будет осуществлён в консоль")
	} else {
		// Закрытие всех открытых файлов в результате настройки логгера
		defer openLogFiles.Close()
	}

	// Создание нового подключения к БД
	db, err := repository.NewPostgresDB(config.PQ)
}
