package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		log.Printf("⏱️ Начало запроса: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		log.Printf("✅ Запрос завершен: %s %s, время выполнения: %v", r.Method, r.URL.Path, duration)
	})
}
