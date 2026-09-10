package middleware

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := rec.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("underlying ResponseWriter does not implement http.Hijacker")
}

func (rec *responseRecorder) Flush() {
	if fl, ok := rec.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

// LoggerMiddleware logs incoming HTTP requests and response latency.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		log.Printf("[HTTP] %s %s | %d | %v | IP: %s", r.Method, r.URL.Path, rec.statusCode, duration, r.RemoteAddr)
	})
}
