package middleware

import (
	"strings"

	"github.com/VictorRibeiroH/ativvo-backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization format (use: Bearer <token>)",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Safely extract user_id from claims. guard against nil or unexpected types to
		// avoid interface{} conversion panics.
		if uidRaw, exists := claims["user_id"]; exists {
			if userID, ok := uidRaw.(string); ok && userID != "" {
				// use the key "user_id" so handlers can read the same name
				c.Locals("user_id", userID)
				return c.Next()
			}
		}
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid token claims",
	})
}
