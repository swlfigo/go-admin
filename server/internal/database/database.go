package database

import (
	"fmt"

	"go-admin/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(c config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
