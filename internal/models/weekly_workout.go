package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WeeklyWorkout struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	
	DayOfWeek int       `json:"day_of_week" gorm:"not null"` // 0-6 (Dom-Sáb)
	Name      string    `json:"name" gorm:"not null"`
	Exercises string    `json:"exercises" gorm:"type:text"` // JSON com lista de exercícios
	IsRest    bool      `json:"is_rest" gorm:"default:false"`
	Completed bool      `json:"completed" gorm:"default:false"`
	
	WeekStart time.Time `json:"week_start" gorm:"not null;index"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (w *WeeklyWorkout) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}
