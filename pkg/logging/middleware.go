package logging

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func Middleware(logger logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			rec := newStatusRecorder(w)

			next.ServeHTTP(rec, req)

			duration := time.Since(start)
			entry := logger.WithFields(logrus.Fields{
				"method":     req.Method,
				"endpoint":   req.URL.Path,
				"status":     rec.status,
				"duration":   duration.String(),
				"latency_ms": float64(duration.Microseconds()) / 1000.0,
			})

			switch {
			case rec.status >= http.StatusInternalServerError:
				entry.Error("request failed")
			case rec.status >= http.StatusBadRequest:
				entry.Warn("request completed")
			default:
				entry.Info("request completed")
			}
		})
	}
}
