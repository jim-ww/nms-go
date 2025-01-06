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
	"github.com/jmoiron/sqlx"

	authHandler "github.com/jim-ww/nms-go/internal/features/auth/handler"
	authMiddleware "github.com/jim-ww/nms-go/internal/features/auth/middleware"
	authService "github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	jwtService "github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	noteHandler "github.com/jim-ww/nms-go/internal/features/note/handler"
	noteService "github.com/jim-ww/nms-go/internal/features/note/services/note"
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

	driverName := "sqlite3"
	db, err := sql.Open(driverName, cfg.StoragePath)
	if err != nil {
		log.Fatal("Failed to initialize sqlite3 db")
	}
	e.Logger.Info("Initialized storage", "storage-path", cfg.StoragePath)

	migrations.MustMigrate(db)
	e.Logger.Info("Database migration completed successfully")

	repo := repository.New(db)
	sqlxDB := sqlx.NewDb(db, driverName)

	jwtService := jwtService.New(cfg.JWTTokenConfig)
	authService := authService.New(jwtService, repo)
	noteService := noteService.New(repo, sqlxDB)

	authHandler := authHandler.New(authService, jwtService)
	noteHandler := noteHandler.New(noteService)

	authMiddleware := authMiddleware.New(jwtService)

	routes.AddRoutes(e, authHandler, authMiddleware, noteHandler)

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
