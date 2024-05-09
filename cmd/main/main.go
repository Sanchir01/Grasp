package main

import (
	"github.com/Sanchir01/Grasp/internal/config"
	"github.com/Sanchir01/Grasp/pkg/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

var (
	development = "development"
	production  = "production"
)

func main() {

	cfg := config.InitConfig()
	log := setupLogger(cfg.Env)
	log.Info("Starting server for", slog.String("port", cfg.HttpServer.Address))

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case development:
		log = setupPrettySlog()

	case production:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
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
