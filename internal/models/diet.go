package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DietPlan - Plano de dieta do usuário
type DietPlan struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      string    `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	
	// Dados do usuário no momento da criação
	Age            int     `json:"age" gorm:"not null"`
	Gender         string  `json:"gender" gorm:"type:varchar(20);not null"`
	Height         float64 `json:"height" gorm:"not null"` // cm
	Weight         float64 `json:"weight" gorm:"not null"` // kg
	ActivityLevel  string  `json:"activity_level" gorm:"type:varchar(50);not null"` // sedentary, lightly_active, moderately_active, very_active, extremely_active
	
	// Cálculos
	TMB            float64 `json:"tmb" gorm:"not null"` // Taxa Metabólica Basal
	TDEE           float64 `json:"tdee" gorm:"not null"` // Total Daily Energy Expenditure (GET)
	
	// Objetivo e estratégia
	Goal           string  `json:"goal" gorm:"type:varchar(20);not null"` // cutting, bulking, maintenance
	TargetCalories float64 `json:"target_calories" gorm:"not null"`
	
	// Macros (em porcentagem)
	ProteinPercent float64 `json:"protein_percent" gorm:"not null"`
	CarbsPercent   float64 `json:"carbs_percent" gorm:"not null"`
	FatPercent     float64 `json:"fat_percent" gorm:"not null"`
	
	// Macros (em gramas)
	ProteinGrams   float64 `json:"protein_grams" gorm:"not null"`
	CarbsGrams     float64 `json:"carbs_grams" gorm:"not null"`
	FatGrams       float64 `json:"fat_grams" gorm:"not null"`
	
	// Refeições
	Meals          []Meal  `json:"meals" gorm:"foreignKey:DietPlanID"`
	
	// Status
	IsActive       bool      `json:"is_active" gorm:"default:true;index"`
	
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (d *DietPlan) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

// Meal - Refeição do plano de dieta
type Meal struct {
	ID         string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DietPlanID string     `json:"diet_plan_id" gorm:"type:uuid;not null;index"`
	DietPlan   DietPlan   `json:"-" gorm:"foreignKey:DietPlanID"`
	
	Name       string     `json:"name" gorm:"not null"` // ex: "Café da manhã", "Almoço", etc
	Time       string     `json:"time"` // ex: "08:00"
	Order      int        `json:"order" gorm:"not null"` // ordem da refeição no dia
	
	MealFoods  []MealFood `json:"meal_foods" gorm:"foreignKey:MealID"`
	
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (m *Meal) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// MealFood - Relação entre refeição e alimento (com quantidade)
type MealFood struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MealID    string    `json:"meal_id" gorm:"type:uuid;not null;index"`
	Meal      Meal      `json:"-" gorm:"foreignKey:MealID"`
	FoodID    string    `json:"food_id" gorm:"type:uuid;not null;index"`
	Food      Food      `json:"food" gorm:"foreignKey:FoodID"`
	
	Quantity  float64   `json:"quantity" gorm:"not null"` // em gramas
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (mf *MealFood) BeforeCreate(tx *gorm.DB) error {
	if mf.ID == "" {
		mf.ID = uuid.New().String()
	}
	return nil
}

// DTOs para criação
type CreateDietPlanInput struct {
	Age            int     `json:"age" validate:"required,min=10,max=120"`
	Gender         string  `json:"gender" validate:"required,oneof=male female"`
	Height         float64 `json:"height" validate:"required,min=50,max=300"`
	Weight         float64 `json:"weight" validate:"required,min=20,max=500"`
	ActivityLevel  string  `json:"activity_level" validate:"required"`
	Goal           string  `json:"goal" validate:"required,oneof=cutting bulking maintenance"`
	ProteinPercent float64 `json:"protein_percent" validate:"required,min=10,max=60"`
	CarbsPercent   float64 `json:"carbs_percent" validate:"required,min=10,max=70"`
	FatPercent     float64 `json:"fat_percent" validate:"required,min=10,max=50"`
}

type CreateMealInput struct {
	Name  string `json:"name" validate:"required"`
	Time  string `json:"time"`
	Order int    `json:"order" validate:"required,min=1"`
}

type AddFoodToMealInput struct {
	FoodID   string  `json:"food_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,min=1"`
}
