package course

import (
	"rpl-service/models"
	"rpl-service/repositories"
)

type CourseRepository struct {
	repositories.Repository[models.Course]
}
