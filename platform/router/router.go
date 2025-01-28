package router

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rpl-service/config"
	"rpl-service/platform/authenticator"
)

// New registers the routes and returns the router.
func New(_ *authenticator.Authenticator, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	config.InitializeRoutes(router, db)

	return router
}
