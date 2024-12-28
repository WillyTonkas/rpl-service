package models

import (
	"gorm.io/gorm"
)

type IsEnrolled struct {
	gorm.Model
	UserID   uint
	CourseID uint
	IsOwner  bool
}

type Profile struct {
	Name       string
	LastName   string
	Email      string
	University string
	Career     string
	Census     int
}
