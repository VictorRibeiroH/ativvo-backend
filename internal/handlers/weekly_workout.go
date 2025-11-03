package handlers

import (
	"encoding/json"
	"time"

	"github.com/VictorRibeiroH/ativvo-backend/internal/database"
	"github.com/VictorRibeiroH/ativvo-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WeeklyWorkoutInput struct {
	DayOfWeek int      `json:"day_of_week" validate:"required,min=0,max=6"`
	Name      string   `json:"name" validate:"required"`
	Exercises []string `json:"exercises"`
	IsRest    bool     `json:"is_rest"`
}

func GetWeeklyWorkouts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	now := time.Now()
	weekStart := getWeekStart(now)

	var workouts []models.WeeklyWorkout
	if err := database.DB.Where("user_id = ? AND week_start = ?", userID, weekStart).
		Order("day_of_week ASC").
		Find(&workouts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch workouts",
		})
	}

	return c.JSON(fiber.Map{
		"workouts":   workouts,
		"week_start": weekStart,
	})
}

func SaveWeeklyWorkouts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var input []WeeklyWorkoutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	now := time.Now()
	weekStart := getWeekStart(now)

	database.DB.Where("user_id = ? AND week_start = ?", userID, weekStart).Delete(&models.WeeklyWorkout{})

	for _, item := range input {
		exercisesJSON, _ := json.Marshal(item.Exercises)

		workout := models.WeeklyWorkout{
			UserID:    userID,
			DayOfWeek: item.DayOfWeek,
			Name:      item.Name,
			Exercises: string(exercisesJSON),
			IsRest:    item.IsRest,
			WeekStart: weekStart,
		}

		if err := database.DB.Create(&workout).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save workout",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Workouts saved successfully",
	})
}

func ToggleWorkoutComplete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	workoutID := c.Params("id")

	var workout models.WeeklyWorkout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userID).First(&workout).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workout not found",
		})
	}

	workout.Completed = !workout.Completed
	if err := database.DB.Save(&workout).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update workout",
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Workout updated",
		"completed": workout.Completed,
	})
}

func GetWeeklyStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	now := time.Now()
	weekStart := getWeekStart(now)

	// Buscar workouts completados
	var completedCount int64
	database.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&models.WeeklyWorkout{}).
		Where("user_id = ? AND week_start = ? AND completed = ? AND is_rest = ?", userID, weekStart, true, false).
		Count(&completedCount)

	// Buscar meta do usuário
	var user models.User
	goal := 3 // meta padrão
	if err := database.DB.Session(&gorm.Session{PrepareStmt: false}).
		Select("weekly_workouts").Where("id = ?", userID).First(&user).Error; err == nil {
		if user.WeeklyWorkouts > 0 {
			goal = user.WeeklyWorkouts
		}
	}

	emoji := getMotivationEmoji(int(completedCount), goal)

	return c.JSON(fiber.Map{
		"completed": completedCount,
		"goal":      goal,
		"emoji":     emoji,
	})
}

func getWeekStart(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	daysToSubtract := weekday - 1
	weekStart := t.AddDate(0, 0, -daysToSubtract)
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
}

func getMotivationEmoji(completed, goal int) string {
	percentage := float64(completed) / float64(goal)

	if percentage >= 1.0 {
		return "🔥🔥🔥"
	} else if percentage >= 0.6 {
		return "🔥"
	} else if percentage >= 0.3 {
		return "👍"
	}
	return "💪"
}
