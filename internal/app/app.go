package app

import (
	"bd_service/internal/app/grpc"
	"bd_service/internal/storage/postgres"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	GRPCServer *grpc.App
}

func NewApp(
	log *slog.Logger,
	grpcPort int,
	db *pgxpool.Pool,
) *App {
	storage := postgres.New(db)

	grpcServer := grpc.New(log, grpcPort, storage)

	return &App{
		GRPCServer: grpcServer,
	}
}
