package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jim-ww/nms-go/internal/config"
	authValidator "github.com/jim-ww/nms-go/internal/features/auth"
	getAuthform "github.com/jim-ww/nms-go/internal/features/auth/handlers/getAuthForm"
	postauth "github.com/jim-ww/nms-go/internal/features/auth/handlers/postAuth"
	authMiddleware "github.com/jim-ww/nms-go/internal/features/auth/middleware"
	authService "github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/features/auth/services/password/bcrypt"
	userRepo "github.com/jim-ww/nms-go/internal/repository"
	sl "github.com/jim-ww/nms-go/internal/utils/loggers/sl"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	cfg := config.MustLoad()
	cfg.JWTTokenConfig.ExpirationDuration = time.Duration(time.Hour * 24 * 7) // TODO read from config

	logger := sl.SetupLogger(cfg.Env)
	logger.Info("Initialized logger", slog.String("env", cfg.Env), slog.String("http-server.adress", cfg.Address))

	conn, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		logger.Info("Failed to initialize sqlite3 db")
		panic(err)
	}
	logger.Info("Initialized storage", slog.Any("storage-path", cfg.StoragePath))

	userRepo := userRepo.New(conn)

	// m, err := migrate.NewWithInstance("")
	// m.Steps(2)

	// userRepo := userRepo.New(logger, db)
	// userRepo.Migrate()

	jwtService := jwt.New(logger, cfg.JWTTokenConfig)
	passwordHasher := bcrypt.New()

	authValidatr := authValidator.New(logger, validator.New())
	authService := authService.New(logger, jwtService, passwordHasher, userRepo, authValidatr)

	authForm := getAuthform.New(logger)
	authHandler := postauth.New(logger, authService, jwtService)

	authMiddleware := authMiddleware.New(logger, jwtService)

	e := echo.New()

	e.Use(echo.WrapMiddleware(authMiddleware.Handler))
	e.Use(middleware.Logger())

	e.Static("/static", "./static")
	e.File("/favicon.ico", "static/favicon.ico")

	e.GET("/login", authForm.Login)
	e.GET("/register", authForm.Register)
	e.GET("/api/login", authHandler.Login)
	e.GET("/api/register", authHandler.Register)

	logger.Info("Starting server...")
	if err = http.ListenAndServe(cfg.HTTPServer.Address, e); err != nil {
		logger.Error("Failed to start http server", slog.String("address", cfg.HTTPServer.Address), sl.Err(err))
	}
}
