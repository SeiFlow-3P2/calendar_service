package main

import (
	"context"
	"log"
	"time"

	"github.com/SeiFlow-3P2/calendar_service/internal/app"
	"github.com/SeiFlow-3P2/calendar_service/internal/configs"
)

func main() {
	// Загружаем переменные окружения
	configs.LoadEnv()

	// Формируем конфигурацию приложения
	cfg := &app.Config{
		Port:         configs.GetEnv("PORT", "9090"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		MongoURI:     configs.GetMongoURI(),
		MongoDB:      configs.GetMongoDB(),
	}

	// Создаём приложение
	app := app.New(cfg)

	// Запускаем приложение
	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}