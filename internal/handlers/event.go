package handlers

import (
	"log"
	"time"

	"github.com/VictorRibeiroH/ativvo-backend/internal/database"
	"github.com/VictorRibeiroH/ativvo-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateEventInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	EventDate   string `json:"event_date" validate:"required"` // formato: 2006-01-02
	EventTime   string `json:"event_time" validate:"required"` // formato: 15:04
}

// CreateEvent cria um novo evento para o usuário autenticado
func CreateEvent(c *fiber.Ctx) error {
	log.Printf("🎯 CreateEvent called")
	
	// Pegar usuário do middleware
	userIDRaw := c.Locals("user_id")
	log.Printf("📝 user_id from Locals: %v (type: %T)", userIDRaw, userIDRaw)
	
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		log.Printf("❌ user_id type assertion failed")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("❌ Failed to parse user_id: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	
	log.Printf("✅ User ID parsed: %s", userID)

	var input CreateEventInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Parse da data
	eventDate, err := time.Parse("2006-01-02", input.EventDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	event := models.Event{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		EventDate:   eventDate,
		EventTime:   input.EventTime,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create event",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event created successfully",
		"event":   event,
	})
}

// GetEvents retorna todos os eventos do usuário autenticado
func GetEvents(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var events []models.Event
	if err := database.DB.Where("user_id = ?", userID).Order("event_date ASC, event_time ASC").Find(&events).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch events",
		})
	}

	return c.JSON(events)
}

// GetEventsByDate retorna eventos de uma data específica
func GetEventsByDate(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	dateStr := c.Query("date") // formato: 2006-01-02
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Date query parameter is required",
		})
	}

	eventDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	var events []models.Event
	if err := database.DB.Where("user_id = ? AND DATE(event_date) = ?", userID, eventDate.Format("2006-01-02")).
		Order("event_time ASC").
		Find(&events).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch events",
		})
	}

	return c.JSON(events)
}

// DeleteEvent deleta um evento
func DeleteEvent(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	eventID := c.Params("id")
	parsedEventID, err := uuid.Parse(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	var event models.Event
	if err := database.DB.Where("id = ? AND user_id = ?", parsedEventID, userID).First(&event).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	if err := database.DB.Delete(&event).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete event",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Event deleted successfully",
	})
}

// UpdateEvent atualiza um evento
func UpdateEvent(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	eventID := c.Params("id")
	parsedEventID, err := uuid.Parse(eventID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	var event models.Event
	if err := database.DB.Where("id = ? AND user_id = ?", parsedEventID, userID).First(&event).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	var input CreateEventInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.EventDate != "" {
		eventDate, err := time.Parse("2006-01-02", input.EventDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Use YYYY-MM-DD",
			})
		}
		event.EventDate = eventDate
	}

	if input.Title != "" {
		event.Title = input.Title
	}
	if input.Description != "" {
		event.Description = input.Description
	}
	if input.EventTime != "" {
		event.EventTime = input.EventTime
	}

	if err := database.DB.Save(&event).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update event",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Event updated successfully",
		"event":   event,
	})
}
