package exercises

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"rpl-service/mappers"
	"rpl-service/models"
	"rpl-service/repositories/course"
	"rpl-service/repositories/exercise"
	"rpl-service/services/users"
)

type ExerciseService struct {
	exerciseRepo exercise.ExerciseRepository
	testRepo     exercise.TestRepository
}

func (s *ExerciseService) FindExercise(exerciseID uuid.UUID, db *gorm.DB) (*models.Exercise, error) {
	return s.exerciseRepo.FindByID(exerciseID, db)
}

func (s *ExerciseService) CreateExercise(db *gorm.DB, exercise models.ExerciseDTO, userID, courseID uuid.UUID) error {
	var courseService = users.CourseService{
		CourseRepository: course.CourseRepository{},
		EnrollRepository: course.EnrollToCourseRepository{},
	}
	if !courseService.IsOwner(db, userID, courseID) {
		return errors.New("this user doesn't have permission to create an exercise")
	}

	testIDs, buildTestErr := s.buildTests(db, exercise)
	if buildTestErr != nil {
		return buildTestErr
	}
	return s.createExercise(db, exercise, testIDs)
}

func (s *ExerciseService) CreateTest(db *gorm.DB, test models.TestDTO) (uuid.UUID, error) {
	currentTest := models.Test{
		Model:      gorm.Model{},
		Name:       test.Name,
		TestScript: test.TestScript,
	}

	createErr := s.testRepo.Create(currentTest, db)
	if createErr != nil {
		return uuid.Nil, createErr
	}
	// TODO: Check if this is the correct way to get the last test ID
	//var currentTestID uuid.UUID
	//db.Model(models.Test{}).Select("ID").Last(&currentTestID)

	return uuid.Nil, nil
}

// SolveExercise exerciseCode should come from the request.
func (s *ExerciseService) SolveExercise(exerciseUUID uuid.UUID, db *gorm.DB, exerciseCode string,
) ([]models.ExerciseResult, error) {
	tests, getTestErr := s.getTestsByExerciseID(exerciseUUID, db)
	if getTestErr != nil {
		return []models.ExerciseResult{}, getTestErr
	}

	currentExercise, FindErr := s.FindExercise(exerciseUUID, db)
	if FindErr != nil {
		return []models.ExerciseResult{}, FindErr
	}

	return s.getTestResults(tests, exerciseCode, currentExercise)
}

// Private methods

func (s *ExerciseService) getTestsByExerciseID(exerciseUUID uuid.UUID, db *gorm.DB) ([]models.Test, error) {
	tests, queryErr := s.testRepo.QueryTests(exerciseUUID, db)
	if queryErr != nil {
		return []models.Test{}, queryErr
	}
	return tests, nil
}

func (s *ExerciseService) createExercise(db *gorm.DB, exercise models.ExerciseDTO, testIDs []string) error {
	return s.exerciseRepo.Create(models.Exercise{
		Model:       gorm.Model{},
		Name:        exercise.Name,
		Description: exercise.Description,
		BaseCode:    exercise.BaseCode,
		TestIDs:     testIDs,
		Points:      exercise.Points,
		UnitNumber:  exercise.UnitNumber,
	}, db)
}

func (s *ExerciseService) buildTests(db *gorm.DB, exercise models.ExerciseDTO) ([]string, error) {
	var testIDs []string
	for _, test := range exercise.TestData {
		ID, createErr := s.CreateTest(db, test)
		if createErr != nil {
			return []string{}, createErr
		}
		testIDs = append(testIDs, ID.String())
	}
	return testIDs, nil
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
