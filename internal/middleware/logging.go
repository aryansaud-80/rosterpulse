package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status  int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}


func (rw *responseWriter) Write(b []byte) (int, error){
	if !rw.wroteHeader {
		rw.ResponseWriter.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(b)
}


func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappped := &responseWriter{
				ResponseWriter: w,
				status:     http.StatusOK,
			}

			next.ServeHTTP(wrappped, r)

			duration := time.Since(start)

			log.Info("HTTP requst processed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrappped.status,
				"duration_ms", duration.Milliseconds(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}
