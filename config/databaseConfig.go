package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"rpl-service/constants"
	"rpl-service/models"
)

func StartDatabase() *gorm.DB {
	// Uncomment these lines when running in local
	// // Retrieve environment variables
	// err := godotenv.Load(".env")
	// if err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//	return nil
	//}

	host := os.Getenv("HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DATABASE_PORT")

	envVariables := []string{host, user, password, dbname, port}

	for _, envVar := range envVariables {
		if envVar == constants.EmptyString {
			log.Fatal("One or more database environment variables are not set")
		}
	}

	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Enable uuid-ossp extension
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatalf("failed to enable uuid-ossp extension: %v", err)
	}

	migrateSchemas(db)

	return db
}

func migrateSchemas(db *gorm.DB) {
	err := db.AutoMigrate(&models.Course{}, &models.Exercise{}, &models.Test{}, &models.IsEnrolled{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
