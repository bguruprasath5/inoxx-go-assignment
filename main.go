package main

import (
	"ionixx/api/controllers"
	"ionixx/api/middlewares"
	"ionixx/api/storage"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

/**
 * Function to setup the user router
 */
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

/**
 * Function to setup the auth router
 */
func setupAuthRouter(r *gin.Engine) {
	authController := &controllers.AuthController{}
	userRouter := r.Group("/auth")
	userRouter.POST("/login", authController.Login)
}

/**
 * Function load environment variables from .env file
 */
func loadEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
}

/**
 * Function to setup the server
 */
func SetupServer() *gin.Engine {
	// load environment variables
	loadEnv()

	// initialize the database
	storage.InitDB()

	// initialize the server
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// setup user router
	setupUserRouter(r)
	// setup auth router
	setupAuthRouter(r)

	// return the server
	return r
}

func main() {
	// initialize the server and start it
	SetupServer().Run(":" + os.Getenv("PORT"))
}
