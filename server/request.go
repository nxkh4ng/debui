package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nxkh4ng/debui/logger"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusWriter{ResponseWriter: w}
		defer func() {
			status := sw.status
			if status == 0 {
				status = http.StatusOK
			}

			logger.Requestf(
				"%s | %s | %s | %s | %s\n",
				time.Now().Format("15:04:05.000"),
				methodColor(r.Method),
				statusColor(status),
				r.URL.Path,
				durationColor(time.Since(start)),
			)
		}()

		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(sw, r)
	})
}

func statusColor(code int) string {
	switch {
	case code >= 500:
		return fmt.Sprintf("\x1b[31m%d\x1b[0m", code)
	case code >= 400:
		return fmt.Sprintf("\x1b[33m%d\x1b[0m", code)
	case code >= 300:
		return fmt.Sprintf("\x1b[36m%d\x1b[0m", code)
	default:
		return fmt.Sprintf("\x1b[32m%d\x1b[0m", code)
	}
}

func durationColor(d time.Duration) string {
	switch {
	case d > 1*time.Second:
		return fmt.Sprintf("\x1b[31m%s\x1b[0m", d.Round(time.Millisecond).String())
	case d > 100*time.Millisecond:
		return fmt.Sprintf("\x1b[33m%s\x1b[0m", d.Round(time.Millisecond).String())
	default:
		return d.Round(time.Millisecond).String()
	}
}

func methodColor(method string) string {
	var color string

	switch method {
	case "GET":
		color = "36"
	case "POST":
		color = "32"
	case "PUT", "PATCH":
		color = "35"
	case "DELETE":
		color = "31"
	case "HEAD", "OPTIONS":
		color = "90"
	default:
		return fmt.Sprintf("%s-6s", method)
	}

	return fmt.Sprintf("\x1b[%sm%-6s\x1b[0m", color, method)
}
