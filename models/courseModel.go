package models

import (
	"github.com/google/uuid"
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
	ID          uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	BaseCode    string      `json:"base-code"`
	Points      int         `json:"points"`
	UnitNumber  int         `json:"unit_number"`
	TestIDs     []uuid.UUID `gorm:"foreignkey:ExerciseId" json:"testIDs"`
}

type Test struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name   string    `json:"name"`
	Input  []string  `json:"input"`
	Output []string  `json:"output"`
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
