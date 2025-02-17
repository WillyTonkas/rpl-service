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

	results, exerciseError := solveExercise(r, db, exerciseID)

	if exerciseError.Status != http.StatusOK {
		http.Error(w, exerciseError.Message, exerciseError.Status)
		return
	}

	byteResults, _ := json.Marshal(results) //nolint:musttag // No need
	_, _ = w.Write(byteResults)
}

func solveExercise(r *http.Request, db *gorm.DB, exerciseID uuid.UUID) ([]models.ExerciseResult, ExerciseError) {
	var response models.SolveExerciseResponse
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		return nil, ExerciseError{Message: "Invalid body format", Status: http.StatusBadRequest}
	}
	if respErr := json.Unmarshal(body, &response); respErr != nil {
		return nil, ExerciseError{Message: "Invalid body format", Status: http.StatusBadRequest}
	}
	results, exerciseError := exerciseService.SolveExercise(exerciseID, db, response.ExerciseCode)
	if exerciseError != nil {
		return nil, ExerciseError{Message: "Error while executing tests", Status: http.StatusInternalServerError}

	}
	return results, ExerciseError{Message: "Solved Successfully", Status: http.StatusOK}
}

func CreateExercise(_ http.ResponseWriter, _ *http.Request, _ *gorm.DB) {
	// TODO
}

func FindExercise(_ http.ResponseWriter, _ *http.Request, _ *gorm.DB) {
	// TODO
}

type ExerciseError struct {
	Message string
	Status  int
}
