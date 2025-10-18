package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=6"` // - esconde do JSON
	Name      string    `json:"name" gorm:"not null" validate:"required,min=2"`
	
	// Dados físicos
	Gender    string  `json:"gender" gorm:"type:varchar(20)"` // male, female, other
	Height    float64 `json:"height"`                         // em cm
	Weight    float64 `json:"weight"`                         // em kg
	BodyFat   float64 `json:"body_fat"`                       // percentual
	
	// Objetivos
	WeeklyWorkouts int    `json:"weekly_workouts"` // treinos por semana
	CardioTime     int    `json:"cardio_time"`     // minutos de cardio
	Goal           string `json:"goal"`            // lose_weight, gain_muscle, maintain
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileInput struct {
	Name           *string  `json:"name,omitempty" validate:"omitempty,min=2"`
	Gender         *string  `json:"gender,omitempty"`
	Height         *float64 `json:"height,omitempty"`
	Weight         *float64 `json:"weight,omitempty"`
	BodyFat        *float64 `json:"body_fat,omitempty"`
	WeeklyWorkouts *int     `json:"weekly_workouts,omitempty"`
	CardioTime     *int     `json:"cardio_time,omitempty"`
	Goal           *string  `json:"goal,omitempty"`
}
