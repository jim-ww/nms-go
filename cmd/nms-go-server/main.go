package main

import (
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jim-ww/nms-go/internal/auth"
	"github.com/jim-ww/nms-go/internal/middleware"
	"github.com/jim-ww/nms-go/internal/repositories/sqlite"
	"github.com/jim-ww/nms-go/pkg/config"
	sl "github.com/jim-ww/nms-go/pkg/utils/loggers/slog"
)

func main() {

	cfg := config.MustLoad()
	cfg.JWTTokenConfig.ExpirationTime = time.Duration(time.Hour * 24 * 7) // TODO

	logger := sl.SetupLogger(cfg.Env)

	logger.Info("Initialized logger", slog.Any("config", cfg)) // TODO remove

	storage := sqlite.NewSqliteStorage(cfg.StoragePath)
	_ = storage

	mux := http.NewServeMux()

	// lh := login.New(authService, logger)

	routes := map[string]http.HandlerFunc{
		"GET /login":    auth.LoginTmpl,
		"GET /register": auth.RegisterTmpl,
		// "POST /api/login":    lh.Login,
		// "POST /api/register": lh.Register,
	}

	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}

	static := http.NewServeMux()
	static.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static"))))
	static.Handle("/favicon.ico", http.FileServer(http.Dir("./web/static")))

	baseTemplate := middleware.WrapHTMXWithBaseTemplate(mux)

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/web/static") || r.URL.Path == "/favicon.ico" {
			static.ServeHTTP(w, r)
		} else {
			baseTemplate.ServeHTTP(w, r)
		}
	})

	log.Fatal(http.ListenAndServe(cfg.HTTPServer.Address, middleware.Logger(mainHandler)))
}
