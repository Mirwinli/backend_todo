package grpc

import (
	"bd_service/internal/server"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
	todo server.Storage
}

func New(log *slog.Logger, port int, todo server.Storage) *App {
	grpcServer := grpc.NewServer()

	server.Register(grpcServer, todo)

	return &App{
		log:  log,
		gRPC: grpcServer,
		port: port,
	}
}

func (a *App) Run() error {
	const op = "grpc.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("starting server")
	if err := a.gRPC.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpc.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server")

	a.gRPC.GracefulStop()
}
