package cat

import (
	"CliniqueVet/config"

	"github.com/go-chi/chi/v5"
)

// Routes définit toutes les routes pour la gestion des chats.
func Routes(configuration *config.Config) *chi.Mux {
	// Création du routeur pour les routes des chats
	router := chi.NewRouter()

	// Initialiser la configuration de gestion des chats avec le repository
	catConfig := New(configuration)

	// Définir les routes
	router.Post("/cats", catConfig.NewCatHandler)                  // Créer un nouveau chat
	router.Get("/cats", catConfig.AllCatHandler)                   // Récupérer tous les chats
	router.Get("/cats/{id:[0-9]+}", catConfig.CatByIDHandler)      // Récupérer un chat par ID
	router.Put("/cats/{id:[0-9]+}", catConfig.ModifyCatHandler)    // Modifier un chat par ID
	router.Delete("/cats/{id:[0-9]+}", catConfig.DeleteCatHandler) // Supprimer un chat par ID

	// Retourner le routeur avec les routes des chats définies
	return router
}
