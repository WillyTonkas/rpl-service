package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"rpl-service/controllers"
)

// Should run the main web application
// This is a mock version of the main application, this is the binary to compile on CD.
func main() {
	// Start the database
	startServer()
}

func startServer() {
	db := startDatabase()

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
	http.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {
		fmt.Println("Hello, World!") // Mock endpoint
	})

	http.HandleFunc(controllers.CourseExistsEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		controllers.CourseExistsEndpoint.HandlerFunction(writer, request, db)
	})

	serverPort := os.Getenv("SERVER_PORT")
	fmt.Println(serverPort)

	if serverPort == "" {
		log.Panic("serverPort environment variable is not set")
	}

	serverError = http.ListenAndServe(":"+serverPort, nil)
	if serverError != nil {
		return
	}
}

func startDatabase() *gorm.DB {
	// Retrieve environment variables
	host := os.Getenv("DATABASE_URL")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DATABASE_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("One or more database environment variables are not set")
	}

	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	return db
}
