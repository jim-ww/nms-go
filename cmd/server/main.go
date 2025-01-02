package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jim-ww/nms-go/internal/config"
	"github.com/jim-ww/nms-go/internal/lib/api"
	"github.com/jim-ww/nms-go/internal/lib/api/routes"

	authHandler "github.com/jim-ww/nms-go/internal/features/auth/handler"
	authMiddleware "github.com/jim-ww/nms-go/internal/features/auth/middleware"
	authService "github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	jwtService "github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	echoLog "github.com/jim-ww/nms-go/internal/lib/loggers/echo"
	"github.com/jim-ww/nms-go/internal/migrations"
	"github.com/jim-ww/nms-go/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()

	e := echo.New()

	e.HTTPErrorHandler = api.GlobalErrorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: cfg.HTTPServer.Timeout,
	}))

	echoLog.SetLevel(e.Logger, cfg.Env)
	e.Logger.Info("Initialized logger", "env", cfg.Env, "http-server.adress", cfg.HTTPServer.Address)

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		log.Fatal("Failed to initialize sqlite3 db")
	}
	e.Logger.Info("Initialized storage", "storage-path", cfg.StoragePath)

	migrations.MustMigrate(db)
	e.Logger.Info("Database migration completed successfully")

	repo := repository.New(db)

	jwtService := jwtService.New(cfg.JWTTokenConfig)
	authService := authService.New(jwtService, repo)
	authHandler := authHandler.New(authService, jwtService)
	authMiddleware := authMiddleware.New(jwtService)

	routes.AddRoutes(e, authHandler, *authMiddleware)

	e.Logger.Info("Starting server...")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		e.Logger.Info("Starting server...")
		if err := e.Start(cfg.HTTPServer.Address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
