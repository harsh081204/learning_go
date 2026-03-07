package config

import "os"

type Config struct {
	AppEnv string
	Port   string
	DBUrl  string
}

func Load() *Config {
	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "dev"),
		Port:   getEnv("PORT", "8080"),
		DBUrl:  getEnv("DB_URL", "postgres://localhost:5432/mydb"),
	}
	return cfg
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
