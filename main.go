package main

import (
	"ionixx/storage"

	"ionixx/controllers"
	"ionixx/middlewares"

	"github.com/gin-gonic/gin"
)

func setupUserRouter(r *gin.Engine) {
	userController := &controllers.UserController{}
	userRouter := r.Group("/user")
	userRouter.GET("/", userController.GetAllUsers)
	userRouter.POST("/", userController.CreateUser)
	userRouter.Use(middlewares.AuthMiddleware)
	{
		userRouter.GET("/:id", userController.GetUserByID)
		userRouter.PUT("/:id", userController.UpdateUserById)
		userRouter.DELETE("/:id", userController.DeleteUserById)
	}
}

func setupAuthRouter(r *gin.Engine) {
	authController := &controllers.AuthController{}
	userRouter := r.Group("/auth")
	userRouter.POST("/login", authController.Login)
}

func main() {
	storage.InitDB()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	setupUserRouter(r)
	setupAuthRouter(r)

	r.Run(":3000")
}
