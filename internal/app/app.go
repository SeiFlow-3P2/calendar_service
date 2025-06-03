package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SeiFlow-3P2/calendar_service/internal/api"
	"github.com/SeiFlow-3P2/calendar_service/internal/configs"
	"github.com/SeiFlow-3P2/calendar_service/internal/interceptor"
	"github.com/SeiFlow-3P2/calendar_service/internal/repository"
	"github.com/SeiFlow-3P2/calendar_service/internal/service"
	pb "github.com/SeiFlow-3P2/calendar_service/pkg/proto/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	MongoURI     string
	MongoDB      string
}

type App struct {
	config      *Config
	mongoClient *mongo.Client
	grpcServer  *grpc.Server
}

func New(cfg *Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Start(ctx context.Context) error {
	// Подключение к MongoDB
	mongoCfg := configs.MongoConfig{
		URI:      a.config.MongoURI,
		Database: a.config.MongoDB,
		Timeout:  10 * time.Second,
	}
	client, err := configs.NewMongoClient(ctx, mongoCfg)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	a.mongoClient = client

	db := client.Database(a.config.MongoDB)

	// Инициализация репозиториев
	eventRepo := repository.NewEventRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Создание индексов
	if err := eventRepo.EnsureIndexes(ctx); err != nil {
		return fmt.Errorf("failed to ensure event indexes: %v", err)
	}
	if err := categoryRepo.EnsureIndexes(ctx); err != nil {
		return fmt.Errorf("failed to ensure category indexes: %v", err)
	}

	// Инициализация сервисов
	eventService := service.NewEventService(eventRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Инициализация хендлеров
	eventHandler := api.NewEventServiceHandler(eventService)
	categoryHandler := api.NewCategoryServiceHandler(categoryService)
	handler := api.NewHandler(eventHandler, categoryHandler)

	// Настройка gRPC-сервера
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthUnaryServerInterceptor()),
	)
	a.grpcServer = grpcServer

	// Регистрация сервиса
	pb.RegisterCalendarServiceServer(grpcServer, handler)

	// Включение reflection для отладки
	reflection.Register(grpcServer)

	// Запуск TCP-сервера
	listener, err := net.Listen("tcp", ":"+a.config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Каналы для обработки ошибок и сигналов
	serverError := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Запуск gRPC-сервера в горутине
	go func() {
		log.Printf("Starting gRPC server on port %s", a.config.Port)
		serverError <- grpcServer.Serve(listener)
	}()

	// Ожидание завершения
	select {
	case err := <-serverError:
		return fmt.Errorf("gRPC server error: %v", err)
	case <-shutdown:
		log.Println("Shutting down gRPC server...")
		// Graceful shutdown
		stopped := make(chan struct{})
		go func() {
			a.grpcServer.GracefulStop()
			close(stopped)
		}()

		// Ожидание завершения graceful shutdown или таймаута
		select {
		case <-stopped:
			log.Println("gRPC server stopped")
		case <-time.After(5 * time.Second):
			log.Println("Graceful shutdown timed out, forcing stop")
			a.grpcServer.Stop()
		}

		// Закрытие соединения с MongoDB
		if err := a.Close(); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}

		return nil
	}
}

func (a *App) Close() error {
	if a.mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.mongoClient.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect MongoDB: %v", err)
		}
	}
	return nil
}