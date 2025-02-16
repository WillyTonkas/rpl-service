package exercise

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"net/http"
	"rpl-service/models"
	"rpl-service/services/exercises"
)

var exerciseService = exercises.ExerciseService{}

func SolveExercise(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// TODO: shorten function
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
	if exercise.Name == "" {
		http.Error(w, "No exercise with that ID", http.StatusNotFound)
		return
	}
	var response models.SolveExerciseResponse
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "Invalid body format", http.StatusBadRequest)
		return
	}
	if respErr := json.Unmarshal(body, &response); respErr != nil {
		http.Error(w, "Invalid body format", http.StatusBadRequest)
		return
	}
	results, exerciseError := exerciseService.SolveExercise(exerciseID, db, response.ExerciseCode)
	if exerciseError != nil {
		http.Error(w, "Error while executing tests", http.StatusInternalServerError)
		return
	}
	byteResults, err := json.Marshal(results) //nolint:musttag // No need
	if err != nil {
		return
	}

	_, writeErr := w.Write(byteResults)
	if writeErr != nil {
		return
	}
}
func CreateExercise(_ http.ResponseWriter, _ *http.Request, _ *gorm.DB) {
	// TODO
}

func FindExercise(_ http.ResponseWriter, _ *http.Request, _ *gorm.DB) {
	// TODO
}
