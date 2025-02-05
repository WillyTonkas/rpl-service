package users

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"rpl-service/models"
)

func EnrollToCourse(db *gorm.DB, userID, courseID uuid.UUID) error {
	// TODO: change the following line
	// if !userExists(db, userID) {
	//	return errors.New("user does not exist")
	//}

	if IsUserInCourse(db, userID, courseID) {
		return errors.New("user is already in course")
	}

	db.Model(models.IsEnrolled{}).Create(models.IsEnrolled{
		Model:    gorm.Model{},
		UserID:   userID,
		CourseID: courseID,
		IsOwner:  false,
	})

	return nil
}

func CreateCourse(db *gorm.DB, userID uuid.UUID, courseName, description string) (models.Course, error) {
	currentCourse := models.Course{
		Model:       gorm.Model{},
		Name:        courseName,
		Description: description,
	}

	if err := db.Model(models.Course{}).Create(&currentCourse).Error; err != nil {
		return models.Course{}, errors.New("error when creating a course")
	}

	isEnrolled := models.IsEnrolled{
		Model:    gorm.Model{},
		UserID:   userID,
		CourseID: currentCourse.ID,
		IsOwner:  true,
	}

	if db.Model(models.IsEnrolled{}).Create(&isEnrolled).Error != nil {
		return models.Course{}, errors.New("error when creating a course")
	}

	return currentCourse, nil
}

func RemoveStudent(db *gorm.DB, userID, courseID, studentID uuid.UUID) error {
	if !IsOwner(db, userID, courseID) {
		return errors.New("this user doesn't have permission to remove any student")
	}

	if !IsUserInCourse(db, studentID, courseID) {
		return errors.New("the user does not exist in the course")
	}

	var student models.IsEnrolled
	db.Model(models.IsEnrolled{}).First(&student, "ID = ?", studentID)
	db.Model(models.IsEnrolled{}).Delete(&student)

	return nil
}

func IsUserInCourse(db *gorm.DB, userID, courseID uuid.UUID) bool {
	if !CourseExists(db, courseID) {
		return false
	}
	return db.Model(models.IsEnrolled{}).Where("UserID = ? AND CourseID = ?", userID, courseID).Error == nil
}

// ------------------------- Private functions -------------------------

func CourseExists(db *gorm.DB, courseID uuid.UUID) bool {
	return db.Model(models.Course{}).Where("ID = ?", courseID).Error == nil
}

// func userExists(db *gorm.DB, id uint) bool {
//	// TODO: use Auth0
//	return true
//}

func IsOwner(db *gorm.DB, userID, courseID uuid.UUID) bool {
	currentUser := models.IsEnrolled{}
	db.Model(models.IsEnrolled{}).Where("UserID = ? AND CourseID = ?", userID, courseID).First(&currentUser)
	return currentUser.IsOwner
}
