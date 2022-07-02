package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		zlog.Debug().Str("handler", r.URL.Path).Str("method", r.Method).Msg("Request")

		lrw := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)

		var logEvent *zerolog.Event
		if lrw.statusCode == http.StatusOK {
			logEvent = zlog.Debug()
		} else if (lrw.statusCode == http.StatusUnauthorized) || (lrw.statusCode == http.StatusBadRequest) {
			logEvent = zlog.Warn()
		} else {
			logEvent = zlog.Error()
		}

		logEvent = logEvent.Str("handler", r.URL.Path).Str("method", r.Method).Int("status", lrw.statusCode)
		elapsed := time.Since(start)
		logEvent.Dur("timing", elapsed).Msg("Response")
	})
}
