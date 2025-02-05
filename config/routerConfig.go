package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"rpl-service/controllers/course"
	"rpl-service/controllers/exercise"
	"rpl-service/models"
	"rpl-service/platform/middleware"
)

func InitializeRoutes(router *gin.Engine, db *gorm.DB) {
	initializeRoutes(router, db, course.Endpoints)
	initializeRoutes(router, db, exercise.Endpoints)
}

// Private methods

func initializeRoutes(router *gin.Engine, db *gorm.DB, endpoints []models.Endpoint) {
	for _, endpoint := range endpoints {
		mapToGinRoute(router, endpoint, db)
	}
}

func mapToGinRoute(router *gin.Engine, endpoint models.Endpoint, db *gorm.DB) {
	methods := map[string]func(string, ...gin.HandlerFunc) gin.IRoutes{
		"GET":    router.GET,
		"POST":   router.POST,
		"PUT":    router.PUT,
		"DELETE": router.DELETE,
	}

	// Validate method exists
	method, exists := methods[endpoint.Method]
	if !exists {
		panic(fmt.Sprintf("Unsupported method: %s", endpoint.Method))
	}

	// Create the base handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		endpoint.HandlerFunction(w, r, db)
	}

	// Apply middleware if the route is protected
	if endpoint.IsProtected {
		method(endpoint.Path, gin.WrapH(middleware.EnsureValidToken()(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler(w, r)
			}))))
	} else {
		// Unprotected route
		method(endpoint.Path, gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}))
	}
}
