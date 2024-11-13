package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func RequestLogger(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s to %s from %s \n", time.Now().Format("2006/01/02 15:04:05"), r.Method, r.RequestURI, r.RemoteAddr)
		f(w, r)
	}
}
