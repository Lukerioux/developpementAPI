package visit

import (
	"CliniqueVet/config"
	"CliniqueVet/database/dbmodel"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config, visitRepository dbmodel.VisitEntryRepository) *chi.Mux {
	visitConfig := New(configuration, visitRepository) // Injecter le repository ici
	router := chi.NewRouter()

	router.Post("/visit", visitConfig.NewVisitHandler)
	router.Get("/historyvisit", visitConfig.VisitHistoryHandler)

	return router
}
