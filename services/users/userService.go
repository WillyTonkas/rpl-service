package users

import (
	"fmt"
	"gorm.io/gorm"
	"rpl-service/models"
)

func CreateCourse(db *gorm.DB, userId uint, courseName, description string) error {
	currentCourse := models.Course{
		Model:       gorm.Model{},
		Name:        courseName,
		Description: description,
	}

	if db.Model(models.Course{}).Create(&currentCourse).Error != nil {
		return fmt.Errorf("Error when creating a course.")
	}

	db.Model(models.IsEnrolled{}).Create(models.IsEnrolled{
		Model:    gorm.Model{},
		UserId:   userId,
		CourseId: currentCourse.ID,
		IsOwner:  true,
	})

	return nil
}

func RemoveStudent(db *gorm.DB, userId, courseId, studentId uint) error {
	if !isOwner(db, userId, courseId) {
		return fmt.Errorf("This user doesn't have permission to create a unit.")
	}

	if !userInCourse(db, studentId, courseId) {
		return fmt.Errorf("The user does not exist in the course.")
	}

	var student models.IsEnrolled
	db.Model(models.IsEnrolled{}).First(&student, "ID = ?", studentId)
	db.Model(models.IsEnrolled{}).Delete(&student)

	return nil
}

func EnrollToCourse(db *gorm.DB, userId, courseId uint) error {
	if userInCourse(db, userId, courseId) {
		return fmt.Errorf("User is already in course.")
	}

	db.Model(models.IsEnrolled{}).Create(models.IsEnrolled{
		Model:    gorm.Model{},
		UserId:   userId,
		CourseId: courseId,
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
