package exercises

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"rpl-service/mappers"
	"rpl-service/models"
	"rpl-service/repositories"
	"rpl-service/services"
	"rpl-service/services/users"
)

type ExerciseService struct {
	services.Service[models.Exercise]
}

func (s *ExerciseService) FindExercise(exerciseID uuid.UUID, db *gorm.DB) (*models.Exercise, error) {
	return s.Repository.FindByID(exerciseID, db)
}

func (s *ExerciseService) CreateExercise(db *gorm.DB, exercise models.ExerciseDTO, userID, courseID uuid.UUID) error {
	var courseService = users.CourseService{
		Service: services.Service[models.Course]{
			Repository: repositories.Repository[models.Course]{},
		},
	}
	if !courseService.IsOwner(db, userID, courseID) {
		return errors.New("this user doesn't have permission to create an exercise")
	}

	testIDs := s.buildTests(db, exercise)

	return s.createExercise(db, exercise, testIDs)
}

func (s *ExerciseService) CreateTest(db *gorm.DB, test models.TestDTO) uuid.UUID {
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
func (s *ExerciseService) SolveExercise(exerciseUUID uuid.UUID, db *gorm.DB, exerciseCode string,
) ([]models.ExerciseResult, error) {
	tests := s.getTestsByExerciseID(exerciseUUID, db)
	exercise, err := s.FindExercise(exerciseUUID, db)

	if err != nil {
		return []models.ExerciseResult{}, err
	}

	return s.getTestResults(tests, exerciseCode, exercise)
}

// Private methods

func (s *ExerciseService) getTestsByExerciseID(exerciseUUID uuid.UUID, db *gorm.DB) []models.Test {
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

func (s *ExerciseService) createExercise(db *gorm.DB, exercise models.ExerciseDTO, testIDs []string) error {
	return s.Repository.Create(models.Exercise{
		Model:       gorm.Model{},
		Name:        exercise.Name,
		Description: exercise.Description,
		BaseCode:    exercise.BaseCode,
		TestIDs:     testIDs,
		Points:      exercise.Points,
		UnitNumber:  exercise.UnitNumber,
	}, db)
}

func (s *ExerciseService) buildTests(db *gorm.DB, exercise models.ExerciseDTO) []string {
	var testIDs []string
	for _, test := range exercise.TestData {
		testIDs = append(testIDs, s.CreateTest(db, test).String())
	}
	return testIDs
}

func (s *ExerciseService) getTestResults(tests []models.Test, exerciseCode string, exercise *models.Exercise,
) ([]models.ExerciseResult, error) {
	var testResults []models.ExerciseResult
	codeRunnerURL := os.Getenv("RUNNER_URL")

	for _, test := range tests {
		reader := mappers.SolveExerciseRequestBody(exerciseCode, test.TestScript)
		response, postErr := http.Post(codeRunnerURL+"/run-tests", "application/json", reader)
		if postErr != nil {
			return []models.ExerciseResult{}, postErr
		}
		testResults = append(testResults, mappers.SolveExerciseResponseToResult(response, exercise.Name, test.Name))

		if closeErr := response.Body.Close(); closeErr != nil {
			return nil, closeErr
		}
	}

	return testResults, nil
}
