package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iamarjun/go-chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	cfg := apiConfig{
		fileserverHits: 0,
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("."))
	fs := cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer))

	rApp := chi.NewRouter()
	rApp.Handle("/", fs)
	rApp.Handle("/assets/", fs)
	r.Mount("/app", rApp)

	rApi := chi.NewRouter()
	rApi.Get("/metrics", cfg.handlerMetrics)
	rApi.Post("/validate_chirp", handlerValidateChirp)
	rApi.Get("/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlerGetChirps(w, r, db)
	})
	rApi.Post("/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlerPostChirps(w, r, db)
	})
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
