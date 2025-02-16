package course

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"rpl-service/config/constants"
	"rpl-service/models"
	"rpl-service/services/users"
)

// Controllers should have all the functions and logic, routers should expose the endpoints.
// This is the controller for the course entity.

// TODO: Create a specialized repository for Courses
var courseService = users.CourseService{}

func Exists(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	courseID := r.PathValue("id") // Get the course ID from the URL
	if courseID == constants.EmptyString {
		http.Error(w, "Course ID is required", http.StatusBadRequest)
		return
	}
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		http.Error(w, "Invalid course ID format", http.StatusBadRequest)
		return
	}

	// Should return whether a user with that ID Exists
	if !courseService.CourseExists(db, courseUUID) {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Course Exists"))
	if err != nil {
		return
	}
}

func Create(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Should create a new course
	var body models.Course
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := uuid.New() // TODO: get the actual userID
	currentCourse, creatingCourseErr := courseService.CreateCourse(db, userID, body.Name, body.Description)
	if creatingCourseErr != nil {
		http.Error(w, "Failed to Create course", http.StatusInternalServerError)
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

// EnrollToCourse TODO: Test this function after implementing auth0.
func EnrollToCourse(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var enrollmentRequest struct {
		UserID   uuid.UUID `json:"UserID"`
		CourseID uuid.UUID `json:"CourseID"`
	}

	if json.NewDecoder(r.Body).Decode(&enrollmentRequest) != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := courseService.EnrollToCourse(db, enrollmentRequest.UserID, enrollmentRequest.CourseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("User enrolled in course successfully"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// StudentExists TODO: Test this function after implementing auth0.
func StudentExists(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var enrollmentRequest struct {
		UserID   uuid.UUID `json:"UserID"`
		CourseID uuid.UUID `json:"CourseID"`
	}

	if json.NewDecoder(r.Body).Decode(&enrollmentRequest) != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !courseService.IsUserInCourse(db, enrollmentRequest.UserID, enrollmentRequest.CourseID) {
		http.Error(w, "User is not enrolled in the course", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("User is enrolled in the course"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// DeleteStudent TODO: Test this function after implementing auth0.
func DeleteStudent(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var deleteRequest struct {
		UserID    uuid.UUID `json:"UserID"`
		CourseID  uuid.UUID `json:"CourseID"`
		StudentID uuid.UUID `json:"StudentID"`
	}

	if json.NewDecoder(r.Body).Decode(&deleteRequest) != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := courseService.RemoveStudent(db, deleteRequest.UserID, deleteRequest.CourseID, deleteRequest.StudentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Student removed from course successfully"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
