package treatment

import (
	"CliniqueVet/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	treatmentConfig := New(configuration)
	router := chi.NewRouter()

	// Route pour récupérer tous les traitements
	router.Get("/treatments", treatmentConfig.AllTreatmentHandler)

	// Route pour créer un traitement
	router.Post("/treatments", treatmentConfig.CreateTreatmentHandler)

	return router
}
