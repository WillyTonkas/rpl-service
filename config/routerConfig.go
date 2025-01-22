package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rpl-service/controllers/course"
	"rpl-service/models"
	"rpl-service/platform/middleware"
)

func InitializeRoutes(router *gin.Engine, db *gorm.DB) {
	for _, endpoint := range course.Endpoints {
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

	if endpoint.IsProtected {
		methods[endpoint.Method](endpoint.Path, middleware.IsAuthenticated, func(ctx *gin.Context) {
			endpoint.HandlerFunction(ctx.Writer, ctx.Request, db)
		})
	} else {
		methods[endpoint.Method](endpoint.Path, func(ctx *gin.Context) {
			endpoint.HandlerFunction(ctx.Writer, ctx.Request, db)
		})
	}
}
