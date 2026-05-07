package exercise

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"net/http"
	"rpl-service/config/constants"
	"rpl-service/models"
	"rpl-service/services/exercises"
)

var exerciseService = exercises.ExerciseService{}

func SolveExercise(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// With exerciseId, I should search for all its tests, run one by one, and then send the result of each one as JSONs
	exerciseID, err := uuid.Parse(r.PathValue("exerciseId"))
	if err != nil {
		http.Error(w, "Invalid exercise ID format", http.StatusBadRequest)
		return
	}

	var exercise models.Exercise
	db.Model(models.Exercise{}).Where(&exercise, "ID = ?", exerciseID)

	// TODO: make a more valid check, exercise.ID is currently a uint, don't know why
	// TODO 2: we may need to remove the gorm.Model and manually declare each UUID as primary key
	if exercise.Name == constants.EmptyString {
		http.Error(w, "No exercise with that ID", http.StatusNotFound)
		return
	}

	results, solveErr := solveExercise(r, db, exerciseID)

	if solveErr.Status != http.StatusOK {
		http.Error(w, solveErr.Message, solveErr.Status)
		return
	}

	byteResults, _ := json.Marshal(results) //nolint:musttag // No need
	_, _ = w.Write(byteResults)
	w.WriteHeader(http.StatusOK)
}

func CreateExercise(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var exerciseDTO models.ExerciseDTO
	_ = json.NewDecoder(r.Body).Decode(&exerciseDTO)
	if createErr := exerciseService.CreateExercise(db, exerciseDTO, uuid.UUID{}, uuid.UUID{}); createErr != nil {
		http.Error(w, "Error while creating exercise", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func FindExercise(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	exerciseID, parseErr := uuid.Parse(r.PathValue("exerciseId"))
	if parseErr != nil {
		http.Error(w, "Invalid exercise ID format", http.StatusBadRequest)
		return
	}

	exercise, exerciseErr := exerciseService.FindExercise(exerciseID, db)
	if exerciseErr != nil {
		http.Error(w, "Error while finding exercise", http.StatusInternalServerError)
		return
	}

	byteExercise, _ := json.Marshal(exercise) //nolint:musttag // No need
	_, _ = w.Write(byteExercise)
	w.WriteHeader(http.StatusCreated)
	return
}

func solveExercise(r *http.Request, db *gorm.DB, exerciseID uuid.UUID) ([]models.ExerciseResult, exerciseError) {
	var response models.SolveExerciseResponse
	body, readErr := io.ReadAll(r.Body)

	if readErr != nil {
		return nil, exerciseError{Message: "Invalid body format", Status: http.StatusBadRequest}
	}
	if respErr := json.Unmarshal(body, &response); respErr != nil {
		return nil, exerciseError{Message: "Invalid body format", Status: http.StatusBadRequest}
	}

	results, exerciseErr := exerciseService.SolveExercise(exerciseID, db, response.ExerciseCode)
	if exerciseErr != nil {
		return nil, exerciseError{Message: "Error while executing tests", Status: http.StatusInternalServerError}
	}
	return results, exerciseError{Message: "Solved Successfully", Status: http.StatusOK}
}

type exerciseError struct {
	Message string
	Status  int
}
