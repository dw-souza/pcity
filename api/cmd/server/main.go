package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dw-souza/pcity/api/internal/config"
	"github.com/dw-souza/pcity/api/internal/googleplaces"
	"github.com/dw-souza/pcity/api/internal/handlers"
	"github.com/dw-souza/pcity/api/internal/repository"
	"github.com/dw-souza/pcity/api/internal/services"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("database ping", "error", err)
		os.Exit(1)
	}

	repo := repository.New(pool)
	placesClient := googleplaces.NewClient(cfg.GooglePlacesAPIKey)
	indexer := services.NewIndexer(repo, placesClient)
	placesSvc := services.NewPlacesService(repo, indexer)
	amenitiesSvc := services.NewAmenitiesService(repo)
	profileSvc := services.NewProfileService(repo)

	api := handlers.NewAPI(
		cfg.SupabaseJWTSecret,
		handlers.NewDevAuthHandler(repo, cfg.SupabaseJWTSecret),
		handlers.NewHealthHandler(),
		handlers.NewCitiesHandler(placesSvc),
		handlers.NewPlacesHandler(placesSvc),
		handlers.NewAmenitiesHandler(amenitiesSvc),
		handlers.NewProfileHandler(profileSvc),
		handlers.NewAdminHandler(indexer),
	)

	handler := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Admin-Key"},
		AllowCredentials: false,
		MaxAge:           300,
	})(api.Router(cfg.AdminKey, cfg.DevAuth))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 300 * time.Second, // indexação em grade pode demorar
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown", "error", err)
	}
}
