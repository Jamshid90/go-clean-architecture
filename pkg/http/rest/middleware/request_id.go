package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)

type ctxKeyRequestID int

const (
	RequestIDKey    ctxKeyRequestID = 0
	RequestIDHeader                 = "X-Request-Id"
)

var (
	prefix string
	reqid  uint64
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			myid := atomic.AddUint64(&reqid, 1)
			id = fmt.Sprintf("%s-%010d", prefix, myid)
		}
		w.Header().Set("X-Request-ID", id)
		ctx = context.WithValue(ctx, RequestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
