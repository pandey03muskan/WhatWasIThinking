package routes

import (
	"TestProject/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.POST("/api/v1/register", controllers.RegisterUser)
	router.GET("/api/v1/user", controllers.GetUser)
	router.POST("/api/v1/login", controllers.LoginUser)
	router.GET("/api/v1/notes/:user_id", controllers.GetNotes)
	router.POST("/api/v1/notes/:user_id", controllers.CreateNote)
}
