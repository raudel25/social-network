package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"social-network-api/internal/models"
)

type Config struct {
	Database struct {
		Host     string `json:"Host"`
		Port     int    `json:"Port"`
		User     string `json:"User"`
		Password string `json:"Password"`
		DBName   string `json:"DBName"`
	} `json:"Database"`
}

var DB *gorm.DB

func ConnectDatabase() {
	dsn := connectionString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database!")
	}

	migrationsDb(db)

	DB = db
}

func migrationsDb(db *gorm.DB) {
	migrateDb(db, &models.User{})

}

func migrateDb(db *gorm.DB, model interface{}) {
	err := db.AutoMigrate(model)
	if err != nil {
		log.Fatalf("Error migrating: %s", err)
	}
}

func connectionString() string {
	jsonFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var config Config
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)

}
