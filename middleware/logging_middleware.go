// middleware/logging_middleware.go
package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware called for:", r.URL.Path) // ← Отладка
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}
