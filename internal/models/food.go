package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Food - Alimento público compartilhado entre todos os usuários
type Food struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string    `json:"name" gorm:"not null;index" validate:"required"`
	Calories     float64   `json:"calories" gorm:"not null" validate:"required,min=0"`
	Protein      float64   `json:"protein" gorm:"not null" validate:"required,min=0"`
	Carbs        float64   `json:"carbs" gorm:"not null" validate:"required,min=0"`
	Fat          float64   `json:"fat" gorm:"not null" validate:"required,min=0"`
	ServingSize  float64   `json:"serving_size" gorm:"default:100"` // em gramas, padrão 100g
	CreatedByID  string    `json:"created_by_id" gorm:"type:uuid"`
	CreatedBy    User      `json:"created_by" gorm:"foreignKey:CreatedByID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (f *Food) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	if f.ServingSize == 0 {
		f.ServingSize = 100
	}
	return nil
}

type CreateFoodInput struct {
	Name        string  `json:"name" validate:"required"`
	Calories    float64 `json:"calories" validate:"required,min=0"`
	Protein     float64 `json:"protein" validate:"required,min=0"`
	Carbs       float64 `json:"carbs" validate:"required,min=0"`
	Fat         float64 `json:"fat" validate:"required,min=0"`
	ServingSize float64 `json:"serving_size"`
}
