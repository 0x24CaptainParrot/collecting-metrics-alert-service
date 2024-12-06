package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoggingHttpMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				status:         0,
				size:           0,
			}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)
			log.Info(
				"HTTP request processed",
				zap.String("method", r.Method),
				zap.String("request uri", r.RequestURI),
				zap.Int("status", lrw.status),
				zap.Int("size", lrw.size),
				zap.Duration("duration", duration),
			)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.status = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
