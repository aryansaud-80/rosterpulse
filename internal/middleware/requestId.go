package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-ID"
)

type contextKey int

const (
	requestCtxIdKey contextKey = iota
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := r.Header.Get(requestIDHeader)

		if strings.TrimSpace(requestID) == "" {
			requestID = uuid.NewString()
		}

		w.Header().Add(requestIDHeader, requestID)
		ctx := context.WithValue(r.Context(), requestCtxIdKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestIdFromContext(ctx context.Context) string {
    requestId := ctx.Value(requestCtxIdKey).(string)
	return requestId
}
