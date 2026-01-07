package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	InternalID int64          `json:"internal_id" db:"internal_id" gorm:"primaryKey"`
	PublicID   uuid.UUID      `json:"public_id" db:"public_id"`
	Name       string         `json:"name" db:"name"`
	Email      string         `json:"email" db:"email" gorm:"unique"`
	Password   string         `json:"password" db:"password" gorm:"column:password"`
	Role       string         `json:"role" db:"role"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
