package handlers

import (
	"strconv"
	"strings"

	"github.com/VictorRibeiroH/ativvo-backend/internal/database"
	"github.com/VictorRibeiroH/ativvo-backend/internal/models"
	"github.com/VictorRibeiroH/ativvo-backend/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var foodValidate = validator.New()

// CreateFood - Criar novo alimento (público)
func CreateFood(c *fiber.Ctx) error {
	var input models.CreateFoodInput
	
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := foodValidate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	food := models.Food{
		Name:        input.Name,
		Calories:    input.Calories,
		Protein:     input.Protein,
		Carbs:       input.Carbs,
		Fat:         input.Fat,
		ServingSize: input.ServingSize,
		CreatedByID: userID.(string),
	}

	if food.ServingSize == 0 {
		food.ServingSize = 100
	}

	if err := database.DB.Create(&food).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create food"})
	}

	return c.Status(fiber.StatusCreated).JSON(food)
}

// GetFoods - Listar todos os alimentos (com busca e paginação)
func GetFoods(c *fiber.Ctx) error {
	var foods []models.Food
	query := database.DB.Session(&gorm.Session{PrepareStmt: false}).Model(&models.Food{})

	// Busca por nome
	if search := c.Query("search"); search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	// Paginação
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Limit(limit).Offset(offset).Preload("CreatedBy").Find(&foods).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch foods"})
	}

	return c.JSON(fiber.Map{
		"foods": foods,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetFood - Buscar alimento por ID
func GetFood(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var food models.Food
	if err := database.DB.Preload("CreatedBy").First(&food, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found"})
	}

	return c.JSON(food)
}

// UpdateFood - Atualizar alimento (apenas o criador)
func UpdateFood(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var food models.Food
	if err := database.DB.First(&food, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found"})
	}

	// Verificar se é o criador
	if food.CreatedByID != userID.(string) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only update foods you created"})
	}

	var input models.CreateFoodInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	food.Name = input.Name
	food.Calories = input.Calories
	food.Protein = input.Protein
	food.Carbs = input.Carbs
	food.Fat = input.Fat
	if input.ServingSize > 0 {
		food.ServingSize = input.ServingSize
	}

	if err := database.DB.Save(&food).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update food"})
	}

	return c.JSON(food)
}

// DeleteFood - Deletar alimento (apenas o criador)
func DeleteFood(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var food models.Food
	if err := database.DB.First(&food, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found"})
	}

	// Verificar se é o criador
	if food.CreatedByID != userID.(string) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only delete foods you created"})
	}

	if err := database.DB.Delete(&food).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete food"})
	}

	return c.JSON(fiber.Map{"message": "Food deleted successfully"})
}

func SearchTACOFoods(c *fiber.Ctx) error {
	query := c.Query("q")
	
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query parameter 'q' is required"})
	}

	tacoFoods, err := services.SearchTACOFoods(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search TACO foods",
			"details": err.Error(),
		})
	}

	results := make([]map[string]interface{}, 0, len(tacoFoods))
	for _, food := range tacoFoods {
		results = append(results, services.ConvertTACOFoodToLocal(food))
	}

	return c.JSON(fiber.Map{
		"foods": results,
		"total": len(results),
		"source": "TACO",
	})
}

func ImportTACOFood(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var input struct {
		TACOID int `json:"taco_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tacoFood, err := services.GetTACOFoodByID(input.TACOID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get food from TACO",
			"details": err.Error(),
		})
	}

	food := models.Food{
		Name:        tacoFood.Description,
		Calories:    tacoFood.Energy,
		Protein:     tacoFood.Protein,
		Carbs:       tacoFood.Carb,
		Fat:         tacoFood.Lipid,
		ServingSize: 100,
		CreatedByID: userID.(string),
	}

	if err := database.DB.Create(&food).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to import food"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Food imported successfully from TACO",
		"food": food,
	})
}
