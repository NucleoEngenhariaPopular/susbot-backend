package database

import (
	"fmt"
	"log"
	"user-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	// Create index for CPF searches
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_cpf ON users(cpf);").Error; err != nil {
		return nil, fmt.Errorf("failed to create CPF index: %v", err)
	}

	DB = db
	log.Println("Successfully connected to database and created indexes")
	return DB, nil
}

func GetDB() *gorm.DB {
	return DB
}

