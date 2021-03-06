package main

import (
	"course/internal/database"
	"course/internal/exercise"
	"course/internal/middleware"
	"course/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	route.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]interface{}{
			"message": "hello world",
		})
	})

	db := database.NewDatabaseConn()
	exerciseService := exercise.NewExerciseService(db)
	userService := user.NewUserService(db)

	// answer
	route.POST("/answer", middleware.Authentication(userService), exerciseService.CreateAnswer)

	// exercises
	route.GET("/exercises/:id", middleware.Authentication(userService), exerciseService.GetExercise)
	route.GET("/exercises/:id/score", middleware.Authentication(userService), exerciseService.GetUserScore)

	route.POST("/excercises", middleware.Authentication(userService), exerciseService.CreateExcercise)

	// question
	route.POST("/question", middleware.Authentication(userService), exerciseService.CreateQuestion)

	// user
	route.POST("/register", userService.Register)
	route.POST("/login", userService.Login)
	route.Run(":8080")
}
