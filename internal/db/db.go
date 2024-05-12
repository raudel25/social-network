package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"social-network-api/internal/models"
)

func ConnectDatabase(config models.Config) *gorm.DB {
	dsn := connectionString(config)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database!")
	}

	migrationsDb(db)

	return db
}

func migrationsDb(db *gorm.DB) {
	migrateDb(db, &models.User{})
	migrateDb(db, &models.Profile{})
	migrateDb(db, &models.Photo{})
	migrateDb(db, &models.Follow{})
	migrateDb(db, &models.Post{})
	migrateDb(db, &models.Message{})
}

func migrateDb(db *gorm.DB, model interface{}) {
	err := db.AutoMigrate(model)
	if err != nil {
		log.Fatalf("Error migrating: %s", err)
	}
}

func connectionString(config models.Config) string {

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)

}
