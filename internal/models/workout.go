package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Workout representa um treino
type Workout struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      string    `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // strength, cardio, flexibility, etc
	Duration    int       `json:"duration"` // em minutos
	Calories    int       `json:"calories"` // estimativa
	
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (w *Workout) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}

type CreateWorkoutInput struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Type        string    `json:"type" validate:"required"`
	Duration    int       `json:"duration" validate:"required,min=1"`
	Calories    int       `json:"calories"`
	Date        time.Time `json:"date" validate:"required"`
}

type UpdateWorkoutInput struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Type        *string    `json:"type,omitempty"`
	Duration    *int       `json:"duration,omitempty"`
	Calories    *int       `json:"calories,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
}
