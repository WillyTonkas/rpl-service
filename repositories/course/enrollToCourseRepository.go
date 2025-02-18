package course

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"rpl-service/models"
	"rpl-service/repositories"
)

type EnrollToCourseRepository struct {
	repositories.Repository[models.IsEnrolled]
}

func (e *EnrollToCourseRepository) RemoveStudent(studentID uuid.UUID, db *gorm.DB) error {
	var student models.IsEnrolled
	db.Model(models.IsEnrolled{}).First(&student, "UserID = ?", studentID)
	return db.Model(models.IsEnrolled{}).Delete(&student).Error
}
