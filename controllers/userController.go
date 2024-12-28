package controllers

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"rpl-service/models"
	"rpl-service/services/users"
)

const BaseURL = "/users"

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

var CourseExistsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            BaseURL + "/course/{id}",
	HandlerFunction: CourseExists,
}
