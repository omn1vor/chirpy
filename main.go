package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/omn1vor/chirpy/internal/dto"
)

func main() {
	godotenv.Load()

	const port = "8080"
	const fileServerPath = "."

	cfg := newApiConfig()
	r := chi.NewRouter()
	corsMux := middlewareCors(r)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	filesHandler := http.FileServer(http.Dir(fileServerPath))
	filesHandler = http.StripPrefix("/app", filesHandler)
	filesHandler = cfg.middlewareMetricsHits(filesHandler)
	r.Handle("/app", filesHandler)
	r.Handle("/app/*", filesHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handleHealth)
	apiRouter.Get("/reset", cfg.resetMetricsHandler)
	apiRouter.Get("/reset", cfg.resetMetricsHandler)
	apiRouter.Get("/chirps", cfg.getChirps)
	apiRouter.Get("/chirps/{id}", cfg.getChirp)
	apiRouter.Post("/chirps", cfg.addChirp)
	apiRouter.Post("/users", cfg.addUser)
	apiRouter.Put("/users", cfg.updateUser)
	apiRouter.Post("/login", cfg.loginUser)
	apiRouter.Post("/refresh", cfg.refreshToken)
	apiRouter.Post("/revoke", cfg.revokeToken)
	r.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.getMetricsHandler)
	r.Mount("/admin", adminRouter)

	log.Printf("Serving files from path %s on port %s\n", fileServerPath, port)
	log.Fatal(server.ListenAndServe())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.Encode(dto.ErrorDto{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't encode chirps to JSON: "+err.Error())
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
