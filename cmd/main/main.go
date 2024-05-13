package main

import (
	"context"
	"github.com/Sanchir01/Grasp/internal/config"
	"github.com/Sanchir01/Grasp/internal/server"
	httpHandlers "github.com/Sanchir01/Grasp/internal/server/http"
	mwlogger "github.com/Sanchir01/Grasp/internal/server/http/middleware/logger"
	"github.com/Sanchir01/Grasp/pkg/lib/logger/handlers/slogpretty"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	development = "development"
	production  = "production"
)

func main() {
	cfg := config.InitConfig()
	log := setupLogger(cfg.Env)
	srv := server.NewHttpServer(cfg)
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer, mwlogger.New(log))
	handlers := httpHandlers.NewChiRouter(r)

	go func() {
		if err := srv.Run(handlers.StartHttpHandlers()); err != nil {
			log.Error("Listen server error", err.Error())
		}
	}()

	log.Info("Listen server staterted", slog.String("port", cfg.HttpServer.Port))

	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
	<-quite

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("Server shutdown error", err.Error())
	}
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
