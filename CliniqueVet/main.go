package main

import (
	"CliniqueVet/config"
	"CliniqueVet/pkg/cat"
	"CliniqueVet/pkg/treatment"
	"CliniqueVet/pkg/visit"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()
	// Routes pour les différentes parties
	router.Mount("/api/v1/cat", cat.Routes(configuration))
	router.Mount("/api/v1/visit", visit.Routes(configuration, configuration.VisitRepository))
	router.Mount("/api/v1/treatment", treatment.Routes(configuration)) // Routes pour les traitements
	return router
}

func main() {
	// Initialisation de la configuration
	configuration, err := config.New()
	if err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}

	// Initialisation des routes
	router := Routes(configuration)

	// Configuration du serveur
	srv := &http.Server{
		Addr:    ":8080", // Le serveur écoute sur le port 8080
		Handler: router,  // Routeur avec les routes définies
	}

	// Lancer le serveur dans un goroutine
	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Attendre un signal d'interruption pour arrêter le serveur
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Arrêt propre du serveur
	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exited")
}
