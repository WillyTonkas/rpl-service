package exercises

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"rpl-service/mappers"
	"rpl-service/models"
	"rpl-service/services/users"
)

func FindExercise(exerciseID uuid.UUID, db *gorm.DB) (models.Exercise, error) {
	var exercise models.Exercise
	db.Model(models.Exercise{}).Where(&exercise, "ID = ?", exerciseID)
	if exercise.Name == "" {
		return models.Exercise{}, errors.New("no exercise with that ID")
	}
	return exercise, nil
}

func CreateExercise(db *gorm.DB, exercise models.ExerciseDTO, userID, courseID uuid.UUID) error {
	if !users.IsOwner(db, userID, courseID) {
		return errors.New("this user doesn't have permission to create an exercise")
	}

	var testIDs []string
	for _, test := range exercise.TestData {
		testIDs = append(testIDs, CreateTest(db, test).String())
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

func CreateTest(db *gorm.DB, test models.TestDTO) uuid.UUID {
	db.Model(models.Test{}).Create(models.Test{
		Model:      gorm.Model{},
		Name:       test.Name,
		TestScript: test.TestScript,
	})

	var currentTestID uuid.UUID
	db.Model(models.Test{}).Select("ID").Last(&currentTestID)

	return currentTestID
}

// SolveExercise exerciseCode should come from the request.
func SolveExercise(exerciseUUID uuid.UUID, db *gorm.DB, exerciseCode string) ([]models.ExerciseResult, error) {
	tests := getTestsByExerciseID(exerciseUUID, db)
	exercise, err := FindExercise(exerciseUUID, db)

	if err != nil {
		return []models.ExerciseResult{}, err
	}

	var testResults []models.ExerciseResult
	codeRunnerURL := os.Getenv("RUNNER_URL")

	for _, test := range tests {
		reader := mappers.SolveExerciseRequestBody(exerciseCode, test.TestScript)
		response, postErr := http.Post(codeRunnerURL+"/run-tests", "application/json", reader)
		if postErr != nil {
			return []models.ExerciseResult{}, postErr
		}
		testResults = append(testResults, mappers.SolveExerciseResponseToResult(response, exercise.Name, test.Name))
		closeErr := response.Body.Close()
		if closeErr != nil {
			return nil, closeErr
		}
	}

	return testResults, nil
}

// Private methods

func getTestsByExerciseID(exerciseUUID uuid.UUID, db *gorm.DB) []models.Test {
	var exercise models.Exercise
	db.First(&exercise, exerciseUUID)

	var tests []models.Test

	// Now I should have the exercise, so I shall get all tests
	for _, testID := range exercise.TestIDs {
		// TODO: see gorm docs and retrieve multiple rows
		var test models.Test
		db.Model(models.Test{}).Where(&test, "ID = ?", testID)
		tests = append(tests, test)
	}

	return tests
}
