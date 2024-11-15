package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/jim-ww/nms-go/internal/config"
	"github.com/jim-ww/nms-go/internal/handlers/login"
	"github.com/jim-ww/nms-go/internal/middleware"
	"github.com/jim-ww/nms-go/internal/repositories/sqlite"
)

func main() {

	cfg := config.MustLoad()
	_ = cfg

	storage := sqlite.NewSqliteStorage(cfg.StoragePath)
	_ = storage

	mux := http.NewServeMux()

	routes := map[string]http.HandlerFunc{
		"GET /login":         login.LoginTmpl,
		"GET /register":      login.RegisterTmpl,
		"POST /api/login":    login.Login,
		"POST /api/register": login.Register,
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
