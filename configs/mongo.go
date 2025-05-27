package configs

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI      string        `yaml:"uri"`      // Пример: "mongodb://localhost:27017"
	Database string        `yaml:"database"` // Имя БД (например, "calendar")
	Timeout  time.Duration `yaml:"timeout"`  // Таймаут подключения
}

// NewMongoClient создаёт и проверяет подключение к MongoDB.
func NewMongoClient(ctx context.Context, cfg MongoConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Проверка подключения
	ctxPing, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	if err = client.Ping(ctxPing, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client, nil
}