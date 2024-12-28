package users

import (
	"errors"
	"gorm.io/gorm"
	"rpl-service/models"
)

func EnrollToCourse(db *gorm.DB, userID, courseID uint) error {
	// TODO: delete the following line
	// if !userExists(db, userID) {
	//	return errors.New("user does not exist")
	//}

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

func CreateCourse(db *gorm.DB, userID uint, courseName, description string) error {
	currentCourse := models.Course{
		Model:       gorm.Model{},
		Name:        courseName,
		Description: description,
	}

	if db.Model(models.Course{}).Create(&currentCourse).Error != nil {
		return errors.New("error when creating a course")
	}

	db.Model(models.IsEnrolled{}).Create(models.IsEnrolled{
		Model:    gorm.Model{},
		UserID:   userID,
		CourseID: currentCourse.ID,
		IsOwner:  true,
	})

	return nil
}

func RemoveStudent(db *gorm.DB, userID, courseID, studentID uint) error {
	if !isOwner(db, userID, courseID) {
		return errors.New("this user doesn't have permission to remove any student")
	}

	if !userInCourse(db, studentID, courseID) {
		return errors.New("the user does not exist in the course")
	}

	var student models.IsEnrolled
	db.Model(models.IsEnrolled{}).First(&student, "ID = ?", studentID)
	db.Model(models.IsEnrolled{}).Delete(&student)

	return nil
}

func CreateExercise(db *gorm.DB, exercise models.ExerciseDTO, userID, courseID uint) error {
	if !isOwner(db, userID, courseID) {
		return errors.New("this user doesn't have permission to create an exercise")
	}

	var testIDs []uint
	for _, test := range exercise.TestData {
		testIDs = append(testIDs, CreateTest(db, test))
	}

	db.Model(models.Exercise{}).Create(models.Exercise{
		Model:       gorm.Model{},
		Name:        exercise.Name,
		Description: exercise.Description,
		BaseCode:    exercise.BaseCode,
		TestIDs:     testIDs,
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

	var currentTestID uint
	db.Model(models.Test{}).Select("ID").Last(&currentTestID)

	return currentTestID
}

// ------------------------- Private functions -------------------------

func isOwner(db *gorm.DB, userID uint, courseID uint) bool {
	currentUser := models.IsEnrolled{}
	db.Model(models.IsEnrolled{}).Where("UserID = ? AND CourseID = ?", userID, courseID).First(&currentUser)
	return currentUser.IsOwner
}

// func userExists(db *gorm.DB, id uint) bool {
//	//TODO: use Auth0
//	return true
//}

func courseExists(db *gorm.DB, courseID uint) bool {
	return db.Model(models.Course{}).Where("ID = ?", courseID).Error != nil
}

func userInCourse(db *gorm.DB, userID, courseID uint) bool {
	if !courseExists(db, courseID) {
		return false
	}
	return db.Model(models.IsEnrolled{}).Where("UserID = ? AND CourseID = ?", userID, courseID).Error != nil
}
