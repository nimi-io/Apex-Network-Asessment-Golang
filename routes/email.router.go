package routes

import (
	"apex-network-assesment/controllers"

	"github.com/gin-gonic/gin"
)

func EmailRoutes(incommingRoutes *gin.RouterGroup) {
	emailGroup := incommingRoutes.Group("/")

	go emailGroup.POST("/send-email", controllers.SendEmail)

}
