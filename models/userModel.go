package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
}

type IsEnrolled struct {
	gorm.Model
	UserId   uint
	CourseId uint
	IsOwner  bool
}

type Course struct {
	gorm.Model
	Name        string
	Description string
}

type Profile struct {
	Name       string
	LastName   string
	Email      string
	University string
	Career     string
	Census     int
}
