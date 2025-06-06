package middlewares

import (
	"log"
	"net/http"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func newStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{w, http.StatusOK}
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		w := newStatusResponseWriter(writer)
		next.ServeHTTP(w, req)
		requestPath := req.URL.Path
		if req.URL.RawQuery != "" {
			requestPath += "?" + req.URL.RawQuery
		}
		log.Printf("%s - %s %s - %d", req.RemoteAddr, req.Method, requestPath, w.statusCode)
	})
}
