package database

import (
	"address-api/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Habilitando trigram
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm;").Error; err != nil {
		return nil, fmt.Errorf("failed to create trigram extension: %v", err)
	}

	// Rodando as migracoes
	err = db.AutoMigrate(&models.UBS{}, &models.Team{}, &models.StreetSegment{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	// Criando trigrams para os nomes daas ruas
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_street_segments_street_name_trgm 
        ON street_segments USING gin (street_name gin_trgm_ops);`).Error; err != nil {
		return nil, fmt.Errorf("failed to create trigram index: %v", err)
	}

	DB = db
	log.Println("Successfully connected to database and created indexes")
	return DB, nil
}

func GetDB() *gorm.DB {
	return DB
}
