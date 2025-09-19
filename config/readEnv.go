package config

import "os"

type PQ struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type Config struct {
	Db       PQ
	LogLevel string
}

func New() *Config {
	return &Config{
		Db: PQ{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "5432"),
			User: getEnv("DB_USER", "postgres"),
			Pass: getEnv("DB_PASS", "postgres"),
			Name: getEnv("DB_NAME", "postgres"),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
