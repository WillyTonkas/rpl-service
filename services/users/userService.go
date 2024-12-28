package users

import (
	"errors"
	"gorm.io/gorm"
	"rpl-service/models"
)

func userExists(db *gorm.DB, id uint) bool {
	return db.Model(models.User{}).Where("ID = ?", id).Error != nil
}

func courseExists(db *gorm.DB, courseID uint) bool {
	return db.Model(models.Course{}).Where("ID = ?", courseID).Error != nil
}

func userInCourse(db *gorm.DB, userID, courseID uint) bool {
	if !courseExists(db, courseID) {
		return false
	}
	return db.Model(models.IsEnrolled{}).Where("UserID = ? AND CourseID = ?", userID, courseID).Error != nil
}

func EnrollToCourse(db *gorm.DB, userID, courseID uint) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}

	if userInCourse(db, userID, courseID) {
		return errors.New("user is already in course")
	}

	db.Model(models.IsEnrolled{}).Create(models.IsEnrolled{
		Model:    gorm.Model{},
		UserID:   userID,
		CourseID: courseID,
	})

	return nil
}
