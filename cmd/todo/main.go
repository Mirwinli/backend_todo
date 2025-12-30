package main

import (
	"bd_service/internal/app"
	"bd_service/internal/config"
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	local = "local"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(local)

	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found, using system environment variables")
	}

	conn, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:4444/postgres?sslmode=disable")
	if err != nil {
		panic("Error connecting to database")
	}

	if err := conn.Ping(context.Background()); err != nil {
		panic("Error pinging database")
	}

	log.Info("starting gRPC server")

	application := app.NewApp(log, cfg.GRPC.Port, conn)

	go application.GRPCServer.Run()
	log.Info("Application started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Info("Stopping gRPC server")
	application.GRPCServer.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
