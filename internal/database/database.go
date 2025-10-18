package database

import (
	"fmt"
	"log"

	"github.com/VictorRibeiroH/ativvo-backend/internal/config"
	"github.com/VictorRibeiroH/ativvo-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error
	
	logLevel := logger.Silent
	if config.AppConfig.IsDevelopment() {
		logLevel = logger.Info
	}

	DB, err = gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL (Supabase)")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Workout{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Database migration completed")
	return nil
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
