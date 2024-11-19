package main

import (
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jim-ww/nms-go/internal/features/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/handlers"
	userRepo "github.com/jim-ww/nms-go/internal/features/user/repository/sqlite"
	"github.com/jim-ww/nms-go/pkg/config"
	"github.com/jim-ww/nms-go/pkg/middleware"
	"github.com/jim-ww/nms-go/pkg/sqlite"
	tmpl "github.com/jim-ww/nms-go/pkg/utils/handlers"
	sl "github.com/jim-ww/nms-go/pkg/utils/loggers/sl"
)

func main() {

	cfg := config.MustLoad()
	cfg.JWTTokenConfig.ExpirationTime = time.Duration(time.Hour * 24 * 7) // TODO

	logger := sl.SetupLogger(cfg.Env)
	logger.Info("Initialized logger", slog.String("env", cfg.Env), slog.String("http-server.adress", cfg.Address))

	storage := sqlite.NewSqliteStorage(cfg.StoragePath)
	logger.Info("Initialized sqlite storage", slog.Any("storage-path", cfg.StoragePath))

	userRepo := userRepo.NewUserRepository(storage)
	logger.Debug("Initialized userRepository")

	// userRepo.Migrate()
	logger.Debug("UserRepo migrated")

	authService := auth.NewAuthService(logger, cfg.JWTTokenConfig, userRepo)
	logger.Debug("Initialized authService")

	tmplHandler := tmpl.NewTmplHandler(logger)
	logger.Debug("Initialized tmplHandler")

	lah := handlers.NewAuthFormHandler(logger, tmplHandler)
	logger.Debug("Initialized authFormHandler")

	lh := handlers.NewAuthHandler(authService, logger, tmplHandler)
	logger.Debug("Initialized authHandler")

	mux := http.NewServeMux()
	routes := map[string]http.HandlerFunc{
		"GET /login":         lah.LoginTmpl,
		"GET /register":      lah.RegisterTmpl,
		"POST /api/login":    lh.Login,
		"POST /api/register": lh.Register,
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

	log.Fatal(http.ListenAndServe(cfg.HTTPServer.Address, middleware.Logger(mainHandler)))
}
