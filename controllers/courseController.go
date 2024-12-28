package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"rpl-service/models"
	"rpl-service/services/users"
)

const BaseURL = "/courses"

func CourseExists(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	courseID := r.URL.Query().Get("id") // Get the course ID from the URL
	if courseID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Should return whether a user with that ID exists
	if !users.CourseExists(db, courseUUID) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateCourse(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Should create a new course
	var body models.Course
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := uuid.New() // TODO: get the actual userID

	currentCourse, creatingCourseErr := users.CreateCourse(db, userID, body.Name, body.Description)
	if creatingCourseErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(currentCourse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

var CourseExistsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            BaseURL + "/course/{id}",
	HandlerFunction: CourseExists,
}

var CreateCourseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/course",
	HandlerFunction: CreateCourse,
}
