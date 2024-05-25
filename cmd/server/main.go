package main

import (
	"log/slog"
	"os"
	"os/signal"
	"ozon/internal/app"
	"ozon/internal/config"
	"ozon/pkg/logger/logsetup"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logsetup.SetupLogger(cfg.Env)
	log.Info("starting application", slog.Any("config", cfg))

	application, err := app.New(log, cfg)
	if err != nil {
		return
	}
	go application.Server.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sgl := <-stop
	log.Info("stopping application", slog.String("signal", sgl.String()))
	application.Server.Stop()
	log.Info("application stopped")
}
