package main

import (
	"TaskSync/internal/service"
	"TaskSync/internal/storage"
	"TaskSync/internal/storage/postgres"
	"TaskSync/internal/transport/http-server/handler"
	"TaskSync/internal/transport/http-server/server"
	"TaskSync/pkg/logger"
	migrations "TaskSync/pkg/migration"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title TaskSync API
// @version 1.0
// @description API Server for Task Tracking
// @host localhost:8080
// @basePath /

func main() {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		slog.Error("invalid .env file", slog.Any("error", err))
	}

	// Настройка логгера
	log := logger.SetupLogger(os.Getenv("ENV"))

	// Конвертация в int
	attempts, err := strconv.Atoi(os.Getenv("DB_ATTEMPTS"))
	if err != nil {
		log.Error("failed conv str to int", slog.Any("error", err))
		panic(err)
	}

	delay, err := strconv.Atoi(os.Getenv("DB_DELAY"))
	if err != nil {
		log.Error("failed conv str to int", slog.Any("error", err))
		panic(err)
	}

	// Настройка подключения к базе данных PostgreSQL
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}, attempts, time.Duration(delay))

	if err != nil {
		log.Error("failed to init PostgresDB", slog.Any("error", err))
		panic(err)
	}

	// Миграции БД
	err = migrations.RunMigrations(db)
	if err != nil {
		log.Error("Failed to create create migrations", slog.Any("error", err))
		panic(err)
	}

	log.Info("Migrations applied successfully!")

	// Инициализация хранилища, сервисов и обработчиков
	repositories := storage.NewStorage(db)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	// Инициализация логгера
	handlers.InitLogger(log)

	// Настройка и запуск сервера
	srv := &server.Server{}
	log.Info("Starting server...")

	go func() {
		if err = srv.Run(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), handlers.InitRouter()); err != nil {
			log.Error("error starting server", slog.Any("error", err))
			panic(err)
		}
	}()

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info("Stopped by Admin", "Signal", sig)
}
