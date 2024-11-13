package main

import (
	"log"
	"net/http"

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

	http.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static"))))
	http.Handle("/favicon.ico", http.FileServer(http.Dir("./web/static")))

	http.HandleFunc("GET /login", middleware.RequestLogger(login.LoginTmpl))
	http.HandleFunc("POST /api/login", middleware.RequestLogger(login.Login))
	http.HandleFunc("GET /register", middleware.RequestLogger(login.RegisterTmpl))

	log.Fatal(http.ListenAndServe(cfg.HTTPServer.Address, nil))
}
