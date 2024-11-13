package middleware

import (
	"fmt"
	"net/http"
)

func LogMiddleware(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s to %s from %s \n", r.Method, r.RequestURI, r.RemoteAddr)
		f(w, r)
	}
}
