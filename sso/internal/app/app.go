package app

import (
	grpcapp "authservis/internal/app/grpc"
	"authservis/internal/services/auth"
	"authservis/internal/storage/postgres"
	"fmt"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	connStr string,
	tokenTTL time.Duration,
) (*App, error) {
	storage, err := postgres.New(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}, nil
}
