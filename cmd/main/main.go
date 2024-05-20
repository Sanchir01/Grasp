package main

import (
	"context"
	"errors"
	"os/signal"
	"syscall"

	"github.com/Sanchir01/Grasp/internal/config"
	"github.com/Sanchir01/Grasp/internal/db/storage"
	"github.com/Sanchir01/Grasp/internal/server"
	httpHandlers "github.com/Sanchir01/Grasp/internal/server/http"
	mwlogger "github.com/Sanchir01/Grasp/internal/server/http/middleware/logger"
	"github.com/Sanchir01/Grasp/pkg/lib/logger/handlers/slogpretty"
	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
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
	srv := server.NewHttpServer(cfg)
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer, mwlogger.New(log))

	db, err := sqlx.Connect("postgres", "user=postgres dbname=golangS sslmode=disable password=sanchirgarik01")
	if err != nil {
		log.Error("sqlx.Connect error", err.Error())
	}
	log.Info("DATABASE CONNECTED", db)
	defer db.Close()

	var (
		myProducts        = storage.NewProductStorage(db)
		categoriesStorage = storage.NewCategoriesStorage(db)
		handlers          = httpHandlers.NewChiRouter(router, cfg, log, myProducts, categoriesStorage)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	go func(ctx context.Context) {
		if err := srv.Run(handlers.StartHttpHandlers()); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Error("Listen server error", slog.String("error", err.Error()))
				return
			}
			log.Error("Listen server error", slog.String("error", err.Error()))
		}
	}(ctx)

	log.Info("Listen server staterted", slog.String("port", cfg.HttpServer.Port))

	<-ctx.Done()
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
