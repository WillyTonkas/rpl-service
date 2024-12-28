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
		IsOwner:  false,
	})

	return nil
}

func CreateExercise(db *gorm.DB, exercise models.ExerciseDTO, userId, courseId uint) error {
	if !isOwner(db, userId, courseId) {
		return fmt.Errorf("This user doesn't have permission to create a unit.")
	}

	testIds := []uint{}
	for _, test := range exercise.TestData {
		testIds = append(testIds, CreateTest(db, test))
	}

	db.Model(models.Exercise{}).Create(models.Exercise{
		Model:       gorm.Model{},
		Name:        exercise.Name,
		Description: exercise.Description,
		BaseCode:    exercise.BaseCode,
		TestIds:     testIds,
		Points:      exercise.Points,
		UnitNumber:  exercise.UnitNumber,
	})

	return nil
}

func CreateTest(db *gorm.DB, test models.TestDTO) uint {
	db.Model(models.Test{}).Create(models.Test{
		Model:  gorm.Model{},
		Name:   test.Name,
		Input:  test.Input,
		Output: test.Output,
	})

	var currentTestId uint
	db.Model(models.Test{}).Select("ID").Last(&currentTestId)

	return currentTestId
}

// ------------------------- Private functions -------------------------

func isOwner(db *gorm.DB, userId uint, courseId uint) bool {
	currentUser := models.IsEnrolled{}
	db.Model(models.IsEnrolled{}).Where("UserId = ? AND CourseId = ?", userId, courseId).First(&currentUser)
	return currentUser.IsOwner
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
