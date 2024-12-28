package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Name        string
	Description string
}

type Exercise struct {
	gorm.Model
	Name        string
	Description string
	BaseCode    string
	Points      int
	UnitNumber  int
	TestIDs     []uint `gorm:"foreignkey:ExerciseId"`
}

type Test struct {
	gorm.Model
	Name   string
	Input  []string
	Output []string
}

type TestDTO struct {
	Name   string
	Input  []string
	Output []string
}

type ExerciseDTO struct {
	Name        string
	Description string
	BaseCode    string
	Points      int
	UnitNumber  int
	TestData    []TestDTO
}
