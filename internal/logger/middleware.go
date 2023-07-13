package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func WithLogging(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		Log.Info("Request processed",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Duration("duration", time.Since(start)),
			zap.Int("status", ww.Status()),
			zap.Int("bytes", ww.BytesWritten()),
		)
	}
	return http.HandlerFunc(logFn)
}
