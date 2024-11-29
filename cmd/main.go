package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jim-ww/nms-go/internal/features/auth/handlers"
	authMiddleware "github.com/jim-ww/nms-go/internal/features/auth/middleware"
	authService "github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/features/auth/services/password/bcrypt"
	noteHandler "github.com/jim-ww/nms-go/internal/features/note/handlers"
	userRepo "github.com/jim-ww/nms-go/internal/features/user/storage/sqlite"
	"github.com/jim-ww/nms-go/pkg/config"
	"github.com/jim-ww/nms-go/pkg/middleware"
	tmpl "github.com/jim-ww/nms-go/pkg/utils/handlers"
	sl "github.com/jim-ww/nms-go/pkg/utils/loggers/sl"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	cfg := config.MustLoad()
	cfg.JWTTokenConfig.ExpirationDuration = time.Duration(time.Hour * 24 * 7) // TODO read from config

	logger := sl.SetupLogger(cfg.Env)
	logger.Info("Initialized logger", slog.String("env", cfg.Env), slog.String("http-server.adress", cfg.Address))

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		logger.Info("Failed to initialize sqlite3 db")
		panic(err)
	}
	logger.Info("Initialized sqlite db", slog.Any("storage-path", cfg.StoragePath))

	userRepo := userRepo.New(logger, db)
	logger.Debug("Initialized userRepository")

	userRepo.Migrate()
	logger.Debug("UserRepo migrated")

	jwtService := jwt.New(logger, cfg.JWTTokenConfig)
	logger.Debug("Initialized jwtService")

	passwordHasher := bcrypt.New()
	authService := authService.New(logger, jwtService, passwordHasher, userRepo)
	logger.Debug("Initialized authService")

	tmplHandler := tmpl.NewTmplHandler(logger)
	logger.Debug("Initialized tmplHandler")

	loginFormHandler := handlers.NewAuthFormHandler(logger, tmplHandler)
	logger.Debug("Initialized loginFormHandler")

	authHandler := handlers.New(logger, authService, jwtService, tmplHandler)
	logger.Debug("Initialized authHandler")

	noteHandler := noteHandler.New(logger, tmplHandler)
	logger.Debug("Initialized noteHandler")

	mux := http.NewServeMux()
	routes := map[string]http.HandlerFunc{
		"GET /login":         loginFormHandler.LoginTmpl,
		"GET /register":      loginFormHandler.RegisterTmpl,
		"POST /api/login":    authHandler.Login,
		"POST /api/register": authHandler.Register,
		"GET /":              noteHandler.Dashboard,
	}

	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
	logger.Debug("Mux http routes set")

	static := http.NewServeMux()
	static.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static"))))
	static.Handle("/favicon.ico", http.FileServer(http.Dir("./web/static")))
	logger.Debug("static mux http routes set")

	htmxMiddleware := middleware.NewHTMXMiddlewareBaseWrapper(logger)
	authMiddleware := authMiddleware.New(logger, jwtService)

	logger.Debug("Initialized htmxMiddlewareBaseWrapper")

	baseTemplate := htmxMiddleware.WrapHTMXWithBaseTemplate(mux)
	logger.Debug("Initialized baseTemplate middleware for mux")

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/web/static") || r.URL.Path == "/favicon.ico" {
			logger.Debug("Got request to /web/static, serving from static mux")
			static.ServeHTTP(w, r)
		} else {
			logger.Debug("Regular request, serving from default mux")
			baseTemplate.ServeHTTP(w, r)
		}
	})

	loggedMux := middleware.Logger(mainHandler)
	protectedMux := authMiddleware.Handler(loggedMux)

	logger.Info("Starting server...")
	if err = http.ListenAndServe(cfg.HTTPServer.Address, protectedMux); err != nil {
		logger.Error("Failed to start http server", slog.String("address", cfg.HTTPServer.Address), sl.Err(err))
	}
}
