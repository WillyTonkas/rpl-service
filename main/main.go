package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	database "rpl-service/config"
	"rpl-service/constants"
	"rpl-service/controllers/course"
)

// Should run the main web application
// This is a mock version of the main application, this is the binary to compile on CD.
func main() {
	// Start the database
	startServer()
}

func startServer() {
	db := database.StartDatabase()
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

	http.HandleFunc(course.ExistsEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		course.ExistsEndpoint.HandlerFunction(writer, request, db)
	})

	http.HandleFunc(course.CreateCourseEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		course.CreateCourseEndpoint.HandlerFunction(writer, request, db)
	})

	http.HandleFunc(course.EnrollToCourseEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		course.EnrollToCourseEndpoint.HandlerFunction(writer, request, db)
	})

	http.HandleFunc(course.StudentExistsEndPoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		course.StudentExistsEndPoint.HandlerFunction(writer, request, db)
	})

	http.HandleFunc(course.DeleteStudentEndpoint.Path, func(writer http.ResponseWriter, request *http.Request) {
		course.DeleteStudentEndpoint.HandlerFunction(writer, request, db)
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
