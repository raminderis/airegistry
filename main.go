package main

import (
	"log"
	"net/http"
	"os"

	"cutmenot.ai/airegistry/registry"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}

	port := getEnv("REGISTRYMGR_LISTENING_PORT", "8081")
	addr := ":" + port

	router := chi.NewRouter()
	router.Post("/registries", registry.AddHandler)
	router.Put("/registries/{id}", registry.UpdateHandler)
	router.Delete("/registries/{id}", registry.DeleteHandler)
	router.Get("/registries/{id}", registry.GetHandler)
	router.Get("/registries/by-name", registry.GetByNameHandler)
	router.Get("/registries", registry.GetAllHandler)
	log.Printf("registries mgr service listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
