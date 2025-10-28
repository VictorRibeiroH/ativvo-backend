package database

import (
	"fmt"

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

	dsn := config.AppConfig.DatabaseURL

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logLevel),
		PrepareStmt: false, // Desabilita prepared statements para evitar duplicação
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Limpar prepared statements antigos
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	_, err = sqlDB.Exec("DEALLOCATE ALL")
	if err != nil {
		return fmt.Errorf("could not deallocate prepared statements: %w", err)
	}

	return nil
}

func Migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Workout{},
		&models.WeeklyWorkout{},
		&models.Event{},
	)
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}


