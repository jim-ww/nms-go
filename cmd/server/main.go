package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/jim-ww/nms-go/internal/config"
	authHandler "github.com/jim-ww/nms-go/internal/features/auth/handler"
	authMiddleware "github.com/jim-ww/nms-go/internal/features/auth/middleware"
	authService "github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	jwtService "github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/middleware/errorhandler"
	userRepo "github.com/jim-ww/nms-go/internal/repository"
	echoLog "github.com/jim-ww/nms-go/internal/utils/loggers/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()

	e := echo.New()

	echoLog.SetLevel(e.Logger, cfg.Env)
	e.Logger.Info("Initialized logger", "env", cfg.Env, "http-server.adress", cfg.HTTPServer.Address)

	conn, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		log.Fatal("Failed to initialize sqlite3 db")
	}
	e.Logger.Info("Initialized storage", "storage-path", cfg.StoragePath)

	// 	sqlite3MigrationDriver, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	// 	if err != nil {
	// 		logger.Error("failed to prepare sqlite3MigrationDriver")
	// 		fmt.Println("asdasdasds")
	// 		panic(err)
	// 	}
	//
	// 	m, err := migrate.NewWithDatabaseInstance(
	// 		"file://migrations",
	// 		"nms",
	// 		sqlite3MigrationDriver,
	// 	)
	// 	err = m.Up() // m.Steps(2)
	// 	if err != nil {
	// 		logger.Error("Failed to migrate database", sl.Err(err))
	// 		fmt.Println("123123")
	// 		panic(err)
	// 	}

	userRepo := userRepo.New(conn)

	jwtService := jwtService.New(cfg.JWTTokenConfig)
	authService := authService.New(jwtService, userRepo)
	authHandler := authHandler.New(authService, jwtService)
	authMiddleware := authMiddleware.New(jwtService)

	e.HTTPErrorHandler = errorhandler.CustomHTTPErrorHandler

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
