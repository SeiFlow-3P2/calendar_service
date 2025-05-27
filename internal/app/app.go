package app

import (
	"context"
	"log"
	"time"

	"github.com/SeiFlow-3P2/calendar_service/configs"
	"github.com/SeiFlow-3P2/calendar_service/internal/repository"
	"github.com/SeiFlow-3P2/calendar_service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	MongoClient  *mongo.Client
	EventService *service.EventService
	// Другие сервисы...
}

func NewApp(ctx context.Context) (*App, error) {
	// Конфигурация MongoDB
	mongoCfg := configs.MongoConfig{
		URI:      "mongodb://localhost:27017",
		Database: "calendar_db",
		Timeout:  5 * time.Second,
	}

	// Подключение к MongoDB
	mongoClient, err := configs.NewMongoClient(ctx, mongoCfg)
	if err != nil {
		return nil, err
	}

	db := mongoClient.Database(mongoCfg.Database)

	// Инициализация репозиториев
	eventRepo := repository.NewEventRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Создаем индексы
	if err := eventRepo.EnsureIndexes(ctx); err != nil {
		return nil, err
	}
	if err := categoryRepo.EnsureIndexes(ctx); err != nil {
		return nil, err
	}

	// Инициализация сервисов
	eventService := service.NewEventService(eventRepo, categoryRepo)

	return &App{
		MongoClient:  mongoClient,
		EventService: eventService,
	}, nil
}

func (a *App) Close() {
	if a.MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.MongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}
}
