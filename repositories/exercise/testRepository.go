package exercise

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"rpl-service/models"
	"rpl-service/repositories"
)

type TestRepository struct {
	repositories.Repository[models.Test]
}

func (t *TestRepository) QueryTests(exerciseID uuid.UUID, db *gorm.DB) ([]models.Test, error) {
	var currentExercise models.Exercise
	var tests []models.Test

	db.First(&currentExercise, exerciseID)
	queryErr := db.Where("ID IN ?", currentExercise.TestIDs).Find(&tests).Error
	return tests, queryErr
}
