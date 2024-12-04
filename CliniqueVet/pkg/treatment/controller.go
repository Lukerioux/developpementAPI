package treatment

import (
	"CliniqueVet/config"
	"CliniqueVet/database/dbmodel"
	"net/http"

	"github.com/go-chi/render"
)

type TreatmentConfig struct {
	*config.Config
}

func New(configuration *config.Config) *TreatmentConfig {
	return &TreatmentConfig{configuration}
}

// Handler pour récupérer tous les traitements
func (config *TreatmentConfig) AllTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	treatments, err := config.TreatmentRepository.FindAll()
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to retrieve treatments"})
		return
	}
	render.JSON(w, r, treatments)
}

// Handler pour créer un nouveau traitement
func (config *TreatmentConfig) CreateTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	var newTreatment dbmodel.TreatmentEntry

	// Décoder le corps de la requête JSON dans un objet TreatmentEntry
	if err := render.DecodeJSON(r.Body, &newTreatment); err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to decode request body"})
		return
	}

	// Créer le traitement dans la base de données
	createdTreatment, err := config.TreatmentRepository.Create(&newTreatment)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to create treatment"})
		return
	}

	// Retourner l'objet créé
	render.JSON(w, r, createdTreatment)
}
