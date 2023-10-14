package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	const port = "8080"
	const fileServerPath = "."

	cfg := apiConfig{}
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
	apiRouter.Post("/validate_chirp", validateChirp)
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
	encoder.Encode(errorDto{
		Error: msg,
	})
}
