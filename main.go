package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	cfg := apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("."))
	fs := cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer))

	rApp := chi.NewRouter()
	rApp.Handle("/", fs)
	rApp.Handle("/assets/", fs)
	r.Mount("/app", rApp)

	rApi := chi.NewRouter()
	rApi.Get("/healthz", handlerReadiness)
	rApi.Get("/metrics", cfg.handlerMetrics)
	r.Mount("/api", rApi)

	rAdmin := chi.NewRouter()
	rAdmin.HandleFunc("/metrics", cfg.handlerMetrics)
	r.Mount("/admin", rAdmin)

	corsMux := middlewareCors(r)

	httpServer := &http.Server{
		ReadTimeout: http.DefaultClient.Timeout,
		Handler:     corsMux,
		Addr:        ":8080",
	}

	httpServer.ListenAndServe()
}
