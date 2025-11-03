package handlers

import (
	"math"

	"github.com/VictorRibeiroH/ativvo-backend/internal/database"
	"github.com/VictorRibeiroH/ativvo-backend/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var dietValidate = validator.New()

// CalculateTMB - Calcula Taxa Metabólica Basal usando fórmula de Mifflin-St Jeor
func CalculateTMB(weight, height float64, age int, gender string) float64 {
	var tmb float64
	
	if gender == "male" {
		// Homens: TMB = 10 * peso(kg) + 6.25 * altura(cm) - 5 * idade + 5
		tmb = (10 * weight) + (6.25 * height) - (5 * float64(age)) + 5
	} else {
		// Mulheres: TMB = 10 * peso(kg) + 6.25 * altura(cm) - 5 * idade - 161
		tmb = (10 * weight) + (6.25 * height) - (5 * float64(age)) - 161
	}
	
	return math.Round(tmb)
}

// CalculateTDEE - Calcula Gasto Energético Total Diário
func CalculateTDEE(tmb float64, activityLevel string) float64 {
	var multiplier float64
	
	switch activityLevel {
	case "sedentary":
		multiplier = 1.2
	case "lightly_active":
		multiplier = 1.375
	case "moderately_active":
		multiplier = 1.55
	case "very_active":
		multiplier = 1.725
	case "extremely_active":
		multiplier = 1.9
	default:
		multiplier = 1.55
	}
	
	return math.Round(tmb * multiplier)
}

// CalculateTargetCalories - Calcula calorias alvo baseado no objetivo
func CalculateTargetCalories(tdee float64, goal string) float64 {
	switch goal {
	case "cutting":
		return math.Round(tdee - 500)
	case "bulking":
		return math.Round(tdee + 500)
	case "maintenance":
		return tdee
	default:
		return tdee
	}
}

// CalculateMacros - Calcula macros em gramas baseado nas porcentagens
func CalculateMacros(calories, proteinPercent, carbsPercent, fatPercent float64) (protein, carbs, fat float64) {
	protein = math.Round((calories * proteinPercent / 100) / 4)
	carbs = math.Round((calories * carbsPercent / 100) / 4)
	fat = math.Round((calories * fatPercent / 100) / 9)
	return
}

// CalculateDietPlan - Endpoint para calcular TMB, TDEE e valores sugeridos
func CalculateDietPlan(c *fiber.Ctx) error {
	var input models.CreateDietPlanInput
	
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := dietValidate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	totalPercent := input.ProteinPercent + input.CarbsPercent + input.FatPercent
	if totalPercent < 99 || totalPercent > 101 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Macros percentages must sum to 100%"})
	}

	tmb := CalculateTMB(input.Weight, input.Height, input.Age, input.Gender)
	tdee := CalculateTDEE(tmb, input.ActivityLevel)
	targetCalories := CalculateTargetCalories(tdee, input.Goal)
	
	proteinGrams, carbsGrams, fatGrams := CalculateMacros(
		targetCalories,
		input.ProteinPercent,
		input.CarbsPercent,
		input.FatPercent,
	)

	return c.JSON(fiber.Map{
		"tmb":             tmb,
		"tdee":            tdee,
		"target_calories": targetCalories,
		"protein_grams":   proteinGrams,
		"carbs_grams":     carbsGrams,
		"fat_grams":       fatGrams,
	})
}

// CreateDietPlan - Criar novo plano de dieta
func CreateDietPlan(c *fiber.Ctx) error {
	var input models.CreateDietPlanInput
	
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := dietValidate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	totalPercent := input.ProteinPercent + input.CarbsPercent + input.FatPercent
	if totalPercent < 99 || totalPercent > 101 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Macros percentages must sum to 100%"})
	}

	database.DB.Model(&models.DietPlan{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false)

	tmb := CalculateTMB(input.Weight, input.Height, input.Age, input.Gender)
	tdee := CalculateTDEE(tmb, input.ActivityLevel)
	targetCalories := CalculateTargetCalories(tdee, input.Goal)
	proteinGrams, carbsGrams, fatGrams := CalculateMacros(
		targetCalories,
		input.ProteinPercent,
		input.CarbsPercent,
		input.FatPercent,
	)

	dietPlan := models.DietPlan{
		UserID:         userID.(string),
		Age:            input.Age,
		Gender:         input.Gender,
		Height:         input.Height,
		Weight:         input.Weight,
		ActivityLevel:  input.ActivityLevel,
		TMB:            tmb,
		TDEE:           tdee,
		Goal:           input.Goal,
		TargetCalories: targetCalories,
		ProteinPercent: input.ProteinPercent,
		CarbsPercent:   input.CarbsPercent,
		FatPercent:     input.FatPercent,
		ProteinGrams:   proteinGrams,
		CarbsGrams:     carbsGrams,
		FatGrams:       fatGrams,
		IsActive:       true,
	}

	if err := database.DB.Create(&dietPlan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create diet plan"})
	}

	return c.Status(fiber.StatusCreated).JSON(dietPlan)
}

// GetActiveDietPlan - Buscar plano de dieta ativo do usuário
func GetActiveDietPlan(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var dietPlan models.DietPlan
	if err := database.DB.
		Preload("Meals.MealFoods.Food").
		Where("user_id = ? AND is_active = ?", userID, true).
		First(&dietPlan).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No active diet plan found"})
	}

	return c.JSON(dietPlan)
}

// GetDietPlans - Listar todos os planos de dieta do usuário
func GetDietPlans(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var dietPlans []models.DietPlan
	if err := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&dietPlans).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch diet plans"})
	}

	return c.JSON(dietPlans)
}

// CreateMeal - Adicionar refeição ao plano de dieta
func CreateMeal(c *fiber.Ctx) error {
	dietPlanID := c.Params("dietPlanId")
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var dietPlan models.DietPlan
	if err := database.DB.First(&dietPlan, "id = ? AND user_id = ?", dietPlanID, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Diet plan not found"})
	}

	var input models.CreateMealInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := dietValidate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	meal := models.Meal{
		DietPlanID: dietPlanID,
		Name:       input.Name,
		Time:       input.Time,
		Order:      input.Order,
	}

	if err := database.DB.Create(&meal).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create meal"})
	}

	return c.Status(fiber.StatusCreated).JSON(meal)
}

// AddFoodToMeal - Adicionar alimento a uma refeição
func AddFoodToMeal(c *fiber.Ctx) error {
	mealID := c.Params("mealId")
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var meal models.Meal
	if err := database.DB.
		Joins("JOIN diet_plans ON diet_plans.id = meals.diet_plan_id").
		Where("meals.id = ? AND diet_plans.user_id = ?", mealID, userID).
		First(&meal).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Meal not found"})
	}

	var input models.AddFoodToMealInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := dietValidate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var food models.Food
	if err := database.DB.First(&food, "id = ?", input.FoodID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found"})
	}

	mealFood := models.MealFood{
		MealID:   mealID,
		FoodID:   input.FoodID,
		Quantity: input.Quantity,
	}

	if err := database.DB.Create(&mealFood).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add food to meal"})
	}

	database.DB.Preload("Food").First(&mealFood, "id = ?", mealFood.ID)

	return c.Status(fiber.StatusCreated).JSON(mealFood)
}

// RemoveFoodFromMeal - Remover alimento de uma refeição
func RemoveFoodFromMeal(c *fiber.Ctx) error {
	mealFoodID := c.Params("mealFoodId")
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var mealFood models.MealFood
	if err := database.DB.
		Joins("JOIN meals ON meals.id = meal_foods.meal_id").
		Joins("JOIN diet_plans ON diet_plans.id = meals.diet_plan_id").
		Where("meal_foods.id = ? AND diet_plans.user_id = ?", mealFoodID, userID).
		First(&mealFood).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found in meal"})
	}

	if err := database.DB.Delete(&mealFood).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove food from meal"})
	}

	return c.JSON(fiber.Map{"message": "Food removed from meal successfully"})
}

// DeleteMeal - Deletar refeição
func DeleteMeal(c *fiber.Ctx) error {
	mealID := c.Params("mealId")
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var meal models.Meal
	if err := database.DB.
		Joins("JOIN diet_plans ON diet_plans.id = meals.diet_plan_id").
		Where("meals.id = ? AND diet_plans.user_id = ?", mealID, userID).
		First(&meal).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Meal not found"})
	}

	database.DB.Where("meal_id = ?", mealID).Delete(&models.MealFood{})

	if err := database.DB.Delete(&meal).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete meal"})
	}

	return c.JSON(fiber.Map{"message": "Meal deleted successfully"})
}
