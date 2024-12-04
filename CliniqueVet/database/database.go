package database

import (
	"CliniqueVet/database/dbmodel"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("CliniqueVet.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&dbmodel.CatEntry{})
	DB.AutoMigrate(&dbmodel.TreatmentEntry{})
	DB.AutoMigrate(&dbmodel.VisitEntry{})
	log.Println("Database connected and migrated")
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&dbmodel.CatEntry{},
		&dbmodel.TreatmentEntry{},
		&dbmodel.VisitEntry{},
	)
	log.Println("Database migrated successfully")
}
