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
	if courseID == constants.EmptyString {
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

// TODO: Test this function after implementing auth0.
func EnrollToCourse(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var enrollmentRequest struct {
		UserID   uuid.UUID `json:"user_id"`
		CourseID uuid.UUID `json:"course_id"`
	}

	if json.NewDecoder(r.Body).Decode(&enrollmentRequest) != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := users.EnrollToCourse(db, enrollmentRequest.UserID, enrollmentRequest.CourseID)
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

// TODO: Test this function after implementing auth0.
func StudentExists(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var enrollmentRequest struct {
		UserID   uuid.UUID `json:"user_id"`
		CourseID uuid.UUID `json:"course_id"`
	}

	if json.NewDecoder(r.Body).Decode(&enrollmentRequest) != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !users.IsUserInCourse(db, enrollmentRequest.UserID, enrollmentRequest.CourseID) {
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

// TODO: Test this function after implementing auth0.
func DeleteStudent(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var deleteRequest struct {
		UserID    uuid.UUID `json:"user_id"`
		CourseID  uuid.UUID `json:"course_id"`
		StudentID uuid.UUID `json:"student_id"`
	}

	if json.NewDecoder(r.Body).Decode(&deleteRequest) != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := users.RemoveStudent(db, deleteRequest.UserID, deleteRequest.CourseID, deleteRequest.StudentID)
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

var EnrollToCourseEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/enroll",
	HandlerFunction: EnrollToCourse,
}

var StudentExistsEndPoint = models.Endpoint{
	Method:          models.POST,
	Path:            BaseURL + "/check-enrollment",
	HandlerFunction: StudentExists,
}

var DeleteStudentEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            BaseURL + "/delete-student",
	HandlerFunction: DeleteStudent,
}
