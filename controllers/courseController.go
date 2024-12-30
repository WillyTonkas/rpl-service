package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"rpl-service/constants"
	"rpl-service/models"
	"rpl-service/services/users"
)

const BaseURL = "/courses"

func CourseExists(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	courseID := r.PathValue("id") // Get the course ID from the URL
	if courseID == constants.EMPTY_STRING {
		http.Error(w, "Course ID is required", http.StatusBadRequest)
		return
	}
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		http.Error(w, "Invalid course ID format", http.StatusBadRequest)
		return
	}

	// Should return whether a user with that ID exists
	if !users.CourseExists(db, courseUUID) { // TODO: change package name
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Course exists"))
	if err != nil {
		return
	}
}

func CreateCourse(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Should create a new course
	var body models.Course
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := uuid.New() // TODO: get the actual userID

	currentCourse, creatingCourseErr := users.CreateCourse(db, userID, body.Name, body.Description)
	if creatingCourseErr != nil {
		http.Error(w, "Failed to create course", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(currentCourse)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

var CourseExistsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            BaseURL + "/course/exists/{id}",
	HandlerFunction: CourseExists,
}

var CreateCourseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/course",
	HandlerFunction: CreateCourse,
}
