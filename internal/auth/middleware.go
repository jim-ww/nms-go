package auth

import "net/http"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mysecrettoken" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// TODO TODO TODO

		next.ServeHTTP(w, r)
	})
}
