package visit

import (
	"CliniqueVet/config"
	"CliniqueVet/database/dbmodel"
	"net/http"

	"github.com/go-chi/render"
)

type VisitConfig struct {
	*config.Config
	VisitRepository dbmodel.VisitEntryRepository // Assurez-vous que le repository est correctement injecté
}

func New(configuration *config.Config, visitRepository dbmodel.VisitEntryRepository) *VisitConfig {
	return &VisitConfig{
		Config:          configuration,
		VisitRepository: visitRepository, // Injectez le repository ici
	}
}

// NewVisitHandler crée une nouvelle entrée de visite
func (config *VisitConfig) NewVisitHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres du body ou des query parameters
	var visitEntry dbmodel.VisitEntry
	if err := render.DecodeJSON(r.Body, &visitEntry); err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to decode request body"})
		return
	}

	// Ajouter la visite à la base de données
	createdVisit, err := config.VisitRepository.Create(&visitEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to create visit"})
		return
	}

	// Retourner l'entrée créée
	render.JSON(w, r, createdVisit)
}

// VisitHistoryHandler récupère toutes les entrées de visites
func (config *VisitConfig) VisitHistoryHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := config.VisitRepository.FindAll()
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve history"})
		return
	}
	render.JSON(w, r, entries)
}
