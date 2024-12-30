package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"rpl-service/constants"
	"rpl-service/controllers"
	"rpl-service/models"
)

// Should run the main web application
// This is a mock version of the main application, this is the binary to compile on CD.
func main() {
	// Start the database
	startServer()
}

func startServer() {
	db := startDatabase()
	if db == nil {
		fmt.Println("Error starting the database")
		return
	}

	s, serverError := db.DB()
	if serverError != nil {
		return
	}

	// Defer its closing
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(s)

	// Here should go the functions for each endpoint

	http.HandleFunc(controllers.CourseExistsEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		controllers.CourseExistsEndpoint.HandlerFunction(writer, request, db)
	})

	http.HandleFunc(controllers.CreateCourseEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		controllers.CreateCourseEndpoint.HandlerFunction(writer, request, db)
	})

	http.HandleFunc(controllers.EnrollToCourseEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		controllers.CreateCourseEndpoint.HandlerFunction(writer, request, db)
	})

	http.HandleFunc(controllers.StudentExistsEndPoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		controllers.StudentExistsEndPoint.HandlerFunction(writer, request, db)
	})

	http.HandleFunc(controllers.DeleteStudentEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		controllers.DeleteStudentEndpoint.HandlerFunction(writer, request, db)
	})

	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == constants.EmptyString {
		log.Panic("serverPort environment variable is not set")
	}

	serverError = http.ListenAndServe(":"+serverPort, nil)
	if serverError != nil {
		return
	}
}

func startDatabase() *gorm.DB {
	// Retrieve environment variables
	err := godotenv.Load(".env")
	if err != nil {
		return nil
	}

	host := os.Getenv("HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DATABASE_PORT")

	if host == constants.EmptyString ||
		user == constants.EmptyString ||
		password == constants.EmptyString ||
		dbname == constants.EmptyString ||
		port == constants.EmptyString {
		log.Fatal("One or more database environment variables are not set")
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
