package dbmodel

import (
	"gorm.io/gorm"
)

type CatEntry struct {
	gorm.Model
	ID     int    `json:"id_chat"`  // Changement de "id" à "ID"
	Nom    string `json:"nom_chat"` // Changement de "nom" à "Nom"
	Age    int    `json:"age_chat"` // Changement de "age" à "Age"
	Taille int    `json:"taille_chat"`
	Poids  int    `json:"poids_chat"`
}

type CatEntryRepository interface {
	Create(entry *CatEntry) (*CatEntry, error)
	FindAll() ([]*CatEntry, error)
	FindByID(id int) (*CatEntry, error) // Méthode ajoutée pour récupérer un chat par ID
	Delete(id int) (int64, error)       // Modification du nom pour correspondre à DeleteByID
	Update(id int, updatedCat *CatEntry) (*CatEntry, error)
	Save(entry *CatEntry) (*CatEntry, error)
}

type catEntryRepository struct {
	db *gorm.DB
}

// NewCatEntryRepository retourne une instance de repository
func NewCatEntryRepository(db *gorm.DB) CatEntryRepository {
	return &catEntryRepository{db: db}
}

func (r *catEntryRepository) Save(entry *CatEntry) (*CatEntry, error) {
	if err := r.db.Save(entry).Error; err != nil {
		return nil, err
	}
	return entry, nil
}

// Create ajoute un nouveau chat dans la base de données
func (r *catEntryRepository) Create(entry *CatEntry) (*CatEntry, error) {
	if err := r.db.Create(entry).Error; err != nil {
		return nil, err
	}
	return entry, nil
}

// FindAll récupère tous les chats
func (r *catEntryRepository) FindAll() ([]*CatEntry, error) {
	var entries []*CatEntry
	if err := r.db.Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

// FindByID récupère un chat par son ID
func (r *catEntryRepository) FindByID(id int) (*CatEntry, error) {
	var entry CatEntry
	if err := r.db.First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

// Delete supprime un chat par son ID
func (r *catEntryRepository) Delete(id int) (int64, error) {
	result := r.db.Delete(&CatEntry{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// Update met à jour les informations d'un chat existant
func (r *catEntryRepository) Update(id int, updatedCat *CatEntry) (*CatEntry, error) {
	var existingCat CatEntry
	if err := r.db.First(&existingCat, id).Error; err != nil {
		return nil, err // Si le chat n'est pas trouvé, retourner l'erreur
	}

	// Mise à jour des champs de l'entité existante avec les nouvelles valeurs
	existingCat.Nom = updatedCat.Nom
	existingCat.Age = updatedCat.Age
	existingCat.Taille = updatedCat.Taille
	existingCat.Poids = updatedCat.Poids

	// Sauvegarder les modifications dans la base de données
	if err := r.db.Save(&existingCat).Error; err != nil {
		return nil, err
	}

	// Retourner le chat mis à jour
	return &existingCat, nil
}
