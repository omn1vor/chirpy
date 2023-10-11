package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	cfg := apiConfig{}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fsHandler := http.FileServer(http.Dir("."))
	strippedHandler := http.StripPrefix("/app", fsHandler)
	metricsOnHandler := cfg.middlewareMetricsHits(strippedHandler)
	mux.Handle("/app/", metricsOnHandler)

	mux.HandleFunc("/healthz", handleHealth)
	mux.HandleFunc("/metrics", cfg.getMetricsHandler)
	mux.HandleFunc("/reset", cfg.resetMetricsHandler)

	log.Println("Starting server on port", port)
	log.Fatal(server.ListenAndServe())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}
