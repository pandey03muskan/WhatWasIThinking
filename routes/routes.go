package routes

import (
	"TestProject/controllers"
	"TestProject/middlewares"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.POST("/api/v1/register", controllers.RegisterUser)
	router.GET("/api/v1/user", middlewares.AuthMiddleware(), controllers.GetUser)
	router.POST("/api/v1/login", controllers.LoginUser)
	router.GET("/api/v1/notes", middlewares.AuthMiddleware(), controllers.GetNotes)
	router.POST("/api/v1/notes", middlewares.AuthMiddleware(), controllers.CreateNote)
	router.PUT("/api/v1/notes/:note_id", middlewares.AuthMiddleware(), controllers.UpdateNote)
	router.DELETE("/api/v1/notes/:note_id", middlewares.AuthMiddleware(), controllers.DeleteNote)
}
