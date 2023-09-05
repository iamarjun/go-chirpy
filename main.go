package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	cfg := apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("."))

	fs := cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer))
	r.Handle("/app", fs)
	r.Handle("/app/assets/", fs)
	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", cfg.handlerMetrics)

	corsMux := middlewareCors(r)

	httpServer := &http.Server{
		ReadTimeout: http.DefaultClient.Timeout,
		Handler:     corsMux,
		Addr:        ":8080",
	}

	httpServer.ListenAndServe()
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	type metrics struct {
		hits int
	}
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Before: %v", cfg.fileserverHits)
		cfg.fileserverHits++
		log.Printf("After: %v", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}
