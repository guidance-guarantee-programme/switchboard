package main

import (
	"fmt"
	"net/http"
)

type LoggingResponseWriter struct {
	status int
	http.ResponseWriter
}

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func NewLoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &LoggingResponseWriter{
			status:         200,
			ResponseWriter: w,
		}
		h.ServeHTTP(lw, r)
		fmt.Printf("%s %s %d\n", r.Method, r.URL.Path, lw.status)
	})
}
