package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"rpl-service/config"
	"rpl-service/constants"
	"rpl-service/platform/router"
)

// Should run the main web application
// This is a mock version of the main application, this is the binary to compile on CD.
func main() {
	// Start the database
	startServer()
}

func startServer() {
	db := config.StartDatabase()
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

	// Initialize authenticator
	//auth, err := authenticator.New()
	//if err != nil {
	//	fmt.Println("Failed to start authenticator")
	//	return
	//}

	// Initialize the ginRouter
	ginRouter := router.New()
	config.InitializeRoutes(ginRouter, db)

	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == constants.EmptyString {
		log.Panic("SERVER_PORT environment variable is not set")
	}

	routerError := ginRouter.Run(":" + serverPort)
	if routerError != nil {
		fmt.Println("Failed to start server")
		return
	}
}
