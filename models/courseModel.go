package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Exercise struct {
	gorm.Model
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	BaseCode    string         `json:"base-code"`
	Points      int            `json:"points"`
	UnitNumber  int            `json:"unit_number"`
	TestIDs     pq.StringArray `json:"testIDs" gorm:"type:text[]"`
}

type Test struct {
	gorm.Model
	ID     uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name   string         `json:"name"`
	Input  pq.StringArray `json:"input" gorm:"type:text[]"`
	Output pq.StringArray `json:"output" gorm:"type:text[]"`
}

type TestDTO struct {
	Name   string   `json:"name"`
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

type ExerciseDTO struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BaseCode    string    `json:"base-code"`
	Points      int       `json:"points"`
	UnitNumber  int       `json:"unit_number"`
	TestData    []TestDTO `json:"test-data"`
}
