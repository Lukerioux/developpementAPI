package config

import (
	"CliniqueVet/database"
	"CliniqueVet/database/dbmodel"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	// Connexion aux repositories
	CatRepository       dbmodel.CatEntryRepository
	VisitRepository     dbmodel.VisitEntryRepository
	TreatmentRepository dbmodel.TreatmentEntryRepository

	// Connexion à la base de données
	DB *gorm.DB
}

func New() (*Config, error) {
	config := Config{}

	// Initialisation de la connexion à la base de données
	databaseSession, err := gorm.Open(sqlite.Open("CliniqueVet.db"), &gorm.Config{})
	if err != nil {
		return &config, err
	}

	// Migration des modèles
	database.Migrate(databaseSession)

	// Initialisation des repositories
	config.CatRepository = dbmodel.NewCatEntryRepository(databaseSession)
	config.VisitRepository = dbmodel.NewVisitEntryRepository(databaseSession)
	config.TreatmentRepository = dbmodel.NewTreatmentEntryRepository(databaseSession)

	// Exposer la connexion à la base de données
	config.DB = databaseSession

	return &config, nil
}
