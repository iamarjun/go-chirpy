package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamarjun/go-chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits   int
	jwtSecret        []byte
	polkaApiKey      string
	accessJwtClaims  jwt.RegisteredClaims
	refreshJwtClaims jwt.RegisteredClaims
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	polkaApiKey := os.Getenv("POLKA_API_KEY")

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		os.Remove("database.json")
	}

	cfg := apiConfig{
		fileserverHits: 0,
		jwtSecret:      []byte(jwtSecret),
		polkaApiKey:    polkaApiKey,
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
	rApi.Get("/healthz", handlerReadiness)
	rApi.Get("/metrics", cfg.handlerMetrics)
	rApi.Post("/validate_chirp", handlerValidateChirp)
	rApi.Get("/chirps", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerGetChirps(w, r, db)
	})
	rApi.Get("/chirps/{chirpId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerGetChirpById(w, r, db)
	})
	rApi.Delete("/chirps/{chirpId}", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerDeleteChirp(w, r, db)
	})
	rApi.Post("/chirps", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerPostChirps(w, r, db)
	})
	rApi.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		handlerPostUsers(w, r, db)
	})
	rApi.Put("/users", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerPutUsers(w, r, db)
	})
	rApi.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerPostLogin(w, r, db)
	})
	rApi.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerRefreshToken(w, r, db)
	})
	rApi.Post("/revoke", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerRevokeToken(w, r, db)
	})
	rApi.Post("/polka/webhooks", func(w http.ResponseWriter, r *http.Request) {
		cfg.handlerPolkaWebhook(w, r, db)
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
