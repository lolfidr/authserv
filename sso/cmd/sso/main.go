package main

import (
	"authservis/internal/app"
	"authservis/internal/config"
	"authservis/internal/lib/logger/handlers/slogpretty"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	application, err := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	if err != nil {
		log.Error("failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("starting application",
		slog.Int("port", cfg.GRPC.Port),
		slog.String("env", cfg.Env),
	)

	go func() {
		if err := application.GRPCSrv.Run(); err != nil {
			log.Error("gRPC server run error", slog.String("error", err.Error()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	sig := <-stop
	log.Info("shutting down application", slog.String("signal", sig.String()))

	// Создаем контекст с таймаутом
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Остановка сервера
	application.GRPCSrv.Stop()
	log.Info("gRPC server stopped gracefully")

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
