package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IsEnrolled struct {
	gorm.Model
	UserID   uuid.UUID
	CourseID uuid.UUID
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
