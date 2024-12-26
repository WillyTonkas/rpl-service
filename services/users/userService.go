package users

import (
	"fmt"
	"gorm.io/gorm"
	"rpl-service/models"
)

func userExists(db *gorm.DB, id uint) bool {
	return db.Model(models.User{}).Where("ID = ?", id).Error != nil
}

func courseExists(db *gorm.DB, courseId uint) bool {
	return db.Model(models.Course{}).Where("ID = ?", courseId).Error != nil
}

func userInCourse(db *gorm.DB, userId, courseId uint) bool {
	if !courseExists(db, courseId) {
		return false
	}
	return db.Model(models.IsEnrolled{}).Where("UserId = ? AND CourseId = ?", userId, courseId).Error != nil
}

func EnrollToCourse(db *gorm.DB, userId, courseId uint) error {
	if !userExists(db, userId) {
		return fmt.Errorf("User does not exist.")
	}

	if userInCourse(db, userId, courseId) {
		return fmt.Errorf("User is already in course.")
	}

	db.Model(models.IsEnrolled{}).Create(models.IsEnrolled{
		Model:    gorm.Model{},
		UserId:   userId,
		CourseId: courseId,
	})

	return nil
}
