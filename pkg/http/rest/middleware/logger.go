package middleware

import (
	"github.com/Jamshid90/go-clean-architecture/pkg/http/rest/response"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			rw := response.NewResponseWriter(w, http.StatusOK)
			next.ServeHTTP(rw, r)
			end := time.Now()

			logger.Info("request",
				zap.String("method", r.Method),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("route", r.URL.Path),
				zap.String("request_id", GetReqID(r.Context())),
				zap.Int("code", rw.StatusCode()),
				zap.Duration("time", end.Sub(start)),
			)
		})
	}
}

