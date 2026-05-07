package exercise

import (
	"rpl-service/models"
	"rpl-service/repositories"
)

type ExerciseRepository struct {
	repositories.Repository[models.Exercise]
}
