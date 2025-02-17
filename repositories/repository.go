package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Interacts with directly database.
type Repository[T any] struct{}

// Returns an instance of the searched entity.
func (r *Repository[T]) FindByID(id uuid.UUID, db *gorm.DB) (*T, error) {
	var entity T
	if err := db.First(&entity, "ID = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Returns whether the searched element exists in DB.
func (r *Repository[T]) Exists(id uuid.UUID, db *gorm.DB) bool {
	var entity T
	return db.First(&entity, "id = ?", id).Error != nil
}

// Returns whether entity could be deleted.
func (r *Repository[T]) Delete(id uuid.UUID, db *gorm.DB) bool {
	var entity T
	return db.Delete(&entity, "id = ?", id).Error == nil
}

func (r *Repository[T]) Create(entity T, db *gorm.DB) error {
	var model T
	createAction := db.Model(&model).Create(entity)
	return createAction.Error
}
