package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jim-ww/nms-go/internal/config"
	authHandler "github.com/jim-ww/nms-go/internal/features/auth/handler"
	authMiddleware "github.com/jim-ww/nms-go/internal/features/auth/middleware"
	authService "github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	jwtService "github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/middleware/errorhandler"
	"github.com/jim-ww/nms-go/internal/migrations"
	userRepo "github.com/jim-ww/nms-go/internal/repository"
	echoLog "github.com/jim-ww/nms-go/internal/utils/loggers/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()

	e := echo.New()

	e.HTTPErrorHandler = errorhandler.CustomHTTPErrorHandler

	echoLog.SetLevel(e.Logger, cfg.Env)
	e.Logger.Info("Initialized logger", "env", cfg.Env, "http-server.adress", cfg.HTTPServer.Address)

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		log.Fatal("Failed to initialize sqlite3 db")
	}
	e.Logger.Info("Initialized storage", "storage-path", cfg.StoragePath)

	migrations.MustMigrate(db)
	e.Logger.Info("Database migration completed successfully")

	userRepo := userRepo.New(db)

	jwtService := jwtService.New(cfg.JWTTokenConfig)
	authService := authService.New(jwtService, userRepo)
	authHandler := authHandler.New(authService, jwtService)
	authMiddleware := authMiddleware.New(jwtService)

	e.Use(authMiddleware.Handler)
	e.Use(middleware.Logger())

	e.Static("/web", "./web")
	e.File("/favicon.ico", "web/favicon.ico")

	e.GET("/login", authHandler.LoginForm)
	e.GET("/register", authHandler.RegisterForm)
	e.GET("/api/login", authHandler.Login)
	e.GET("/api/register", authHandler.Register)

	e.Logger.Info("Starting server...")
	if err = http.ListenAndServe(cfg.HTTPServer.Address, e); err != nil {
		e.Logger.Error("Failed to start http server", "address", cfg.HTTPServer.Address, "error", err)
	}
}
