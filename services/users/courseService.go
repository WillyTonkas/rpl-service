package users

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"rpl-service/models"
	"rpl-service/repositories/course"
)

type CourseService struct {
	courseRepository course.CourseRepository
	enrollRepository course.EnrollToCourseRepository
}

func (c *CourseService) EnrollToCourse(db *gorm.DB, userID, courseID uuid.UUID) error {
	if c.IsUserInCourse(db, userID, courseID) {
		return errors.New("user is already in course")
	}

	isEnrolled := models.IsEnrolled{
		Model:    gorm.Model{},
		UserID:   userID,
		CourseID: courseID,
		IsOwner:  false,
	}

	if createErr := c.enrollRepository.Create(isEnrolled, db); createErr != nil {
		return errors.New("error when enrolling user to course")
	}

	return nil
}

func (c *CourseService) CreateCourse(db *gorm.DB, userID uuid.UUID, courseName,
	description string) (models.Course, error) {
	currentCourse := models.Course{
		Model:       gorm.Model{},
		Name:        courseName,
		Description: description,
	}

	if createCourseErr := c.courseRepository.Create(currentCourse, db); createCourseErr != nil {
		return models.Course{}, errors.New("error when creating a course")
	}

	isEnrolled := models.IsEnrolled{
		Model:    gorm.Model{},
		UserID:   userID,
		CourseID: currentCourse.ID,
		IsOwner:  true,
	}

	if createEnrollErr := c.enrollRepository.Create(isEnrolled, db); createEnrollErr != nil {
		return models.Course{}, errors.New("error when creating a course")
	}

	return currentCourse, nil
}

func (c *CourseService) RemoveStudent(db *gorm.DB, userID, courseID, studentID uuid.UUID) error {
	if !c.IsOwner(db, userID, courseID) {
		return errors.New("this user doesn't have permission to remove any student")
	}

	if !c.IsUserInCourse(db, studentID, courseID) {
		return errors.New("the user does not exist in the course")
	}

	if removeErr := c.enrollRepository.RemoveStudent(studentID, db); removeErr != nil {
		return errors.New("error when removing student from course")
	}

	return nil
}

func (c *CourseService) IsUserInCourse(db *gorm.DB, userID, courseID uuid.UUID) bool {
	if !c.CourseExists(db, courseID) {
		return false
	}
	return c.courseRepository.FindUserInCourse(userID, courseID, db)
}

func (c *CourseService) CourseExists(db *gorm.DB, courseID uuid.UUID) bool {
	return c.courseRepository.Exists(courseID, db)
}

func (c *CourseService) IsOwner(db *gorm.DB, userID, courseID uuid.UUID) bool {
	return c.courseRepository.FindUserInCourse(userID, courseID, db)
}
