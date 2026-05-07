package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
