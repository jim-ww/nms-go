package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/jim-ww/nms-go/internal/config"
	getAuthform "github.com/jim-ww/nms-go/internal/features/auth/handlers/getAuthForm"
	sl "github.com/jim-ww/nms-go/internal/utils/loggers/sl"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	cfg := config.MustLoad()
	cfg.JWTTokenConfig.ExpirationDuration = time.Duration(time.Hour * 24 * 7) // TODO read from config

	logger := sl.SetupLogger(cfg.Env)
	logger.Info("Initialized logger", slog.String("env", cfg.Env), slog.String("http-server.adress", cfg.Address))

	_, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		logger.Info("Failed to initialize sqlite3 db")
		panic(err)
	}
	logger.Info("Initialized sqlite db", slog.Any("storage-path", cfg.StoragePath))

	// m, err := migrate.NewWithInstance("")
	// m.Steps(2)

	// userRepo := userRepo.New(logger, db)
	// userRepo.Migrate()

	// jwtService := jwt.New(logger, cfg.JWTTokenConfig)
	// passwordHasher := bcrypt.New()
	// authService := authService.New(logger, jwtService, passwordHasher, userRepo)
	// tmplHandler := tmpl.NewTmplHandler(logger)
	// loginFormHandler := handlers.NewAuthFormHandler(logger, tmplHandler)
	// authHandler := handlers.New(logger, authService, jwtService, tmplHandler)
	// noteHandler := noteHandler.New(logger, tmplHandler)

	// mux := http.NewServeMux()
	// routes := map[string]http.HandlerFunc{
	// 	"GET /login":         loginFormHandler.LoginTmpl,
	// 	"GET /register":      loginFormHandler.RegisterTmpl,
	// 	"POST /api/login":    authHandler.Login,
	// 	"POST /api/register": authHandler.Register,
	// 	"GET /":              noteHandler.Dashboard,
	// }

	// 	for path, handler := range routes {
	// 		mux.HandleFunc(path, handler)
	// 	}
	//
	// 	static := http.NewServeMux()
	// 	static.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// 	static.Handle("/favicon.ico", http.FileServer(http.Dir("./static")))
	//
	// 	htmxMiddleware := middleware.NewHTMXMiddlewareBaseWrapper(logger)
	// 	authMiddleware := authMiddleware.New(logger, jwtService)
	//
	// 	baseTemplate := htmxMiddleware.WrapHTMXWithBaseTemplate(mux)
	//
	// 	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		if strings.HasPrefix(r.URL.Path, "/static") || r.URL.Path == "/favicon.ico" {
	// 			static.ServeHTTP(w, r)
	// 		} else {
	// 			baseTemplate.ServeHTTP(w, r)
	// 		}
	// 	})
	//
	// 	loggedMux := middleware.Logger(mainHandler)
	// 	protectedMux := authMiddleware.Handler(loggedMux)

	// validatr := validator.New()
	e := echo.New()

	e.Static("/static", "./static")
	e.File("/favicon.ico", "static/favicon.ico")

	authForm := getAuthform.New(logger)
	e.GET("/login", authForm.Login)
	e.GET("/register", authForm.Register)

	logger.Info("Starting server...")
	if err = http.ListenAndServe(cfg.HTTPServer.Address, e); err != nil {
		logger.Error("Failed to start http server", slog.String("address", cfg.HTTPServer.Address), sl.Err(err))
	}
}
