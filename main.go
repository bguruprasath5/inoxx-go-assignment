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
	// create the user controller
	userController := &controllers.UserController{}

	// create the user router group
	userRouter := r.Group("/user")

	// user routes
	userRouter.GET("/", userController.GetAllUsers)
	userRouter.POST("/", userController.CreateUser)

	//use auth middleware for all user routes below
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
	// create the auth controller
	authController := &controllers.AuthController{}
	// create the auth router group
	authRouter := r.Group("/auth")
	// auth routes
	authRouter.POST("/login", authController.Login)
}

/**
 * Function load environment variables from .env file
 */
func loadEnv() {
	// load .env file
	err := godotenv.Load(".env")

	// if there is an error, panic and exit
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
