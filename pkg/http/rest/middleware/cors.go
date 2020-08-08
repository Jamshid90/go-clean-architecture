package middleware

import (
	"net/http"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("origin")
		if len(origin) == 0 {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Headers",
			"accept, accept-encoding, authorization, content-type, dnt, origin, user-agent, x-csrftoken, x-requested-with, access-control-allow-origin")
		w.Header().Set("Access-Control-Allow-Methods",
			"DELETE, GET, OPTIONS, PATCH, POST, PUT, PROPFIND, DELETE, HEAD, COPY, MKCOL")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
