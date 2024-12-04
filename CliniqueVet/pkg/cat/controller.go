package cat

import (
	"CliniqueVet/config"
	"CliniqueVet/database/dbmodel"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type CatConfig struct {
	CatRepository dbmodel.CatEntryRepository
}

func New(configuration *config.Config) *CatConfig {
	// Initialiser CatConfig avec le repository déjà configuré
	return &CatConfig{CatRepository: configuration.CatRepository}
}

// NewCatHandler crée un nouveau chat dans la base de données
func (config *CatConfig) NewCatHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de la requête GET
	nomStr := r.URL.Query().Get("name")
	ageStr := r.URL.Query().Get("age")
	tailleStr := r.URL.Query().Get("taille")
	poidsStr := r.URL.Query().Get("poids")

	// Convertir les valeurs en types appropriés
	age, err := strconv.Atoi(ageStr)
	if err != nil || age < 0 {
		http.Error(w, "Invalid age parameter", http.StatusBadRequest)
		return
	}

	taille, err := strconv.Atoi(tailleStr)
	if err != nil || taille < 0 {
		http.Error(w, "Invalid taille parameter", http.StatusBadRequest)
		return
	}

	poids, err := strconv.Atoi(poidsStr)
	if err != nil || poids < 0 {
		http.Error(w, "Invalid poids parameter", http.StatusBadRequest)
		return
	}

	// Créer une nouvelle entrée de chat
	cat := dbmodel.CatEntry{
		Nom:    nomStr,
		Age:    age,
		Taille: taille,
		Poids:  poids,
	}

	// Ajouter l'entrée dans la base de données via le repository
	if _, err := config.CatRepository.Create(&cat); err != nil {
		http.Error(w, fmt.Sprintf("Error saving cat to database: %v", err), http.StatusInternalServerError)
		return
	}

	// Répondre avec le chat ajouté, converti en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

// CatByIDHandler récupère un chat par son ID
func (config *CatConfig) CatByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Utiliser chi.URLParam pour récupérer l'ID dans l'URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	// Chercher le chat dans la base de données
	cat, err := config.CatRepository.FindByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Cat not found")
		return
	}

	// Répondre avec le chat trouvé, converti en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

// AllCatHandler récupère tous les chats
func (config *CatConfig) AllCatHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer tous les chats
	cats, err := config.CatRepository.FindAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving cats: %v", err), http.StatusInternalServerError)
		return
	}

	// Répondre avec la liste des chats, convertie en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}

// DeleteCatHandler supprime un chat par son ID
func (config *CatConfig) DeleteCatHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID du chat depuis l'URL (avec chi)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	// Appel de la méthode Delete de CatRepository qui retourne 2 valeurs
	rowsAffected, err := config.CatRepository.Delete(id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting cat: %v", err))
		return
	}

	if rowsAffected == 0 {
		sendError(w, http.StatusNotFound, "Cat not found")
		return
	}

	// Répondre avec un message de succès
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cat deleted successfully"})
}
func (config *CatConfig) ModifyCatHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis les paramètres de requête (exemple : ?id=2)
	idStr := r.URL.Query().Get("id") // On récupère l'ID depuis la requête
	if idStr == "" {
		sendError(w, http.StatusBadRequest, "ID parameter is required")
		return
	}

	// Convertir l'ID en entier
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	// Récupérer les nouveaux paramètres du chat
	nomStr := r.URL.Query().Get("name")
	ageStr := r.URL.Query().Get("age")
	tailleStr := r.URL.Query().Get("taille")
	poidsStr := r.URL.Query().Get("poids")

	// Convertir les valeurs en types appropriés
	age, err := parseIntegerParam(ageStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	taille, err := parseIntegerParam(tailleStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	poids, err := parseIntegerParam(poidsStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Chercher le chat existant dans la base de données
	cat, err := config.CatRepository.FindByID(id)
	if err != nil {
		sendError(w, http.StatusNotFound, "Cat not found")
		return
	}

	// Modifier les informations du chat
	cat.Nom = nomStr
	cat.Age = age
	cat.Taille = taille
	cat.Poids = poids

	// Sauvegarder les modifications dans la base de données
	updatedCat, err := config.CatRepository.Save(cat)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating cat: %v", err))
		return
	}

	// Répondre avec le chat mis à jour en format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCat)
}

func parseIntegerParam(param string) (int, error) {
	value, err := strconv.Atoi(param)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("Invalid %s parameter", param)
	}
	return value, nil
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
