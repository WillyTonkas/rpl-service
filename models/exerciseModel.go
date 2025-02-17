package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SolveExerciseResponse struct {
	ExerciseCode string `json:"code"`
}

type ExerciseResult struct {
	ExerciseName string
	TestName     string
	TestPassed   bool
	Stdout       string
	Stderr       string
}
type Exercise struct {
	gorm.Model
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	BaseCode    string         `json:"base-code"`
	Points      int            `json:"points"`
	UnitNumber  int            `json:"unit_number"`
	TestIDs     pq.StringArray `json:"testIDs" gorm:"type:text[]"`
}

type Test struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"`
	Name       string    `json:"name"`
	TestScript string    `json:"testScript"`
}

// TestDTO Should remove DTOs?
type TestDTO struct {
	Name       string `json:"name"`
	TestScript string `json:"testScript"`
}

// ExerciseDTO Should remove DTOs?
type ExerciseDTO struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BaseCode    string    `json:"base-code"`
	Points      int       `json:"points"`
	UnitNumber  int       `json:"unit_number"`
	TestData    []TestDTO `json:"test-data"`
}
