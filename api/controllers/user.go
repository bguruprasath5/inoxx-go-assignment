package controllers

import (
	"fmt"
	"ionixx/api/models"
	"ionixx/api/storage"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"ionixx/api/response"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

type UserController struct{}

/**
 * Function to get all users
 * @api {get} /users/ Get all users
 * @apiSuccessExample {json} Success-Response:
 *    HTTP/1.1 200 OK
 *   {
 *    	"success": true,
 *   	"message": "Users fetched successfully!",
 *  	"data": [
 *   		{
 *   			"id": 1,
 *  			"user_name": "test",
 *				"full_name": "Test User",
 * 				"dob": "2020-01-01T00:00:00Z",
 * 				"linkedin_url": "",
 * 				"created_at": "2020-01-01T00:00:00Z",
 * 				"updated_at": "2020-01-01T00:00:00Z"
 *  		},
 *		]
 *   }
 * @apiErrorExample {json} Error-Response:
 *   HTTP/1.1 400 Bad Request
 *  {
 *  	"success": false,
 * 		"message": "Error fetching users!",
 * 		"data": null
 * 	}
 */
func (u *UserController) GetAllUsers(c *gin.Context) {
	var list []models.User
	// Find all users in the database
	result := storage.DB.Find(&list)

	if result.Error != nil {
		// If there is an error, we return a bad request
		response.ErrorJSON(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	// Return a success response with the users
	response.SuccessJSON(c, http.StatusOK, "Users fetched successfully!", list)
}

/**
 * Function to get user by id
 * @api {get} /users/:id Get user by id
 * @apiSuccessExample {json} Success-Response:
 *   HTTP/1.1 200 OK
 *  {
 * 		"success": true,
 * 		"message": "User fetched successfully!",
 * 		"data": {
 * 			"id": 1,
 * 			"user_name": "test",
 * 			"full_name": "Test User",
 * 			"dob": "2020-01-01T00:00:00Z",
 * 			"linkedin_url": "",
 * 			"created_at": "2020-01-01T00:00:00Z",
 * 			"updated_at": "2020-01-01T00:00:00Z"
 * 		}
 * 	}
 * @apiErrorExample {json} Error-Response:
 *  HTTP/1.1 400 Bad Request
 * 	{
 * 		"success": false,
 * 		"message": "Error fetching user!",
 * 		"data": null
 * 	}
 */
func (u *UserController) GetUserByID(c *gin.Context) {
	var user models.User
	// Get the first user with the id
	result := storage.DB.Where("id = ?", c.Param("id")).First(&user)

	if result.Error != nil {
		// If there is an error, we return a bad request
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}
	// Return a success response with the user
	response.SuccessJSON(c, http.StatusOK, "User fetched successfully!", user)
}

type CreateUserRequest struct {
	UserName    string    `json:"user_name" binding:"required,min=3"`
	Password    string    `json:"password" binding:"required,min=6"`
	FullName    string    `json:"full_name" binding:"required"`
	Dob         time.Time `json:"dob" binding:"required"`
	LinkedinURL string    `json:"linkedin_url" binding:"required"`
}

func (c *CreateUserRequest) toUser() (*models.User, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(c.Password), 14)
	if err != nil {
		return nil, err
	}
	return &models.User{
		UserName:    c.UserName,
		Password:    string(hashedPasswordBytes),
		FullName:    c.FullName,
		Dob:         &c.Dob,
		LinkedinURL: c.LinkedinURL,
	}, nil
}

/**
 * Function to create user
 * @api {post} /users/ Create user
 * @apiSuccessExample {json} Success-Response:
 *  HTTP/1.1 200 OK
 * {
 * 	"success": true,
 * 	"message": "User created successfully!",
 * 	"data": {
 * 		"id": 1,
 * 		"user_name": "test",
 * 		"full_name": "Test User",
 * 		"dob": "2020-01-01T00:00:00Z",
 * 		"linkedin_url": "",
 * 		"created_at": "2020-01-01T00:00:00Z",
 * 		"updated_at": "2020-01-01T00:00:00Z"
 * 	}
 * }
 * @apiErrorExample {json} Error-Response:
 * HTTP/1.1 400 Bad Request
 * {
 * 	"success": false,
 * 	"message": "Error creating user!",
 * 	"data": null
 * }
 */
func (u *UserController) CreateUser(c *gin.Context) {
	var userData CreateUserRequest
	// Bind the request body to the user data
	if err := c.ShouldBindJSON(&userData); err != nil {
		// Return the first validation error
		for _, fieldErr := range err.(validator.ValidationErrors) {
			response.ErrorJSON(c, http.StatusBadRequest, fmt.Sprintf("Validation Error : %s %s", fieldErr.Field(), fieldErr.Tag()))
			return
		}
		return
	}

	// Create the user model
	user, userErr := userData.toUser()
	// If there is an error, we return a bad request
	if userErr != nil {
		response.ErrorJSON(c, http.StatusBadRequest, userErr.Error())
		return
	}

	// Save the user to the database
	result := storage.DB.Create(user)

	// If there is an error, we return a bad request
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	// Return a success response
	response.SuccessJSON(c, http.StatusOK, "User created successfully!", nil)

}

type UpdateUserRequest struct {
	FullName    string    `json:"full_name"`
	Dob         time.Time `json:"dob"`
	LinkedinURL string    `json:"linkedin_url"`
}

/**
 * Function to update user by id
 * @api {put} /users/:id Update user by id
 * @apiSuccessExample {json} Success-Response:
 * HTTP/1.1 200 OK
 * {
 * 	"success": true,
 * 	"message": "User updated successfully!",
 * 	"data": {
 * 		"id": 1,
 * 		"user_name": "test",
 * 		"full_name": "Test User",
 * 		"dob": "2020-01-01T00:00:00Z",
 * 		"linkedin_url": "",
 * 		"created_at": "2020-01-01T00:00:00Z",
 * 		"updated_at": "2020-01-01T00:00:00Z"
 * 	}
 * }
 * @apiErrorExample {json} Error-Response:
 * HTTP/1.1 400 Bad Request
 * {
 * 	"success": false,
 * 	"message": "Error updating user!",
 * 	"data": null
 * }
 *
 */
func (u *UserController) UpdateUserById(c *gin.Context) {
	var userData UpdateUserRequest
	// Bind the request body to the user data
	if err := c.BindJSON(&userData); err != nil {
		response.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	// Get the user by id
	user := models.User{}
	result := storage.DB.Where("id = ?", c.Param("id")).First(&user)

	// If user not found, return a bad request
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	// If FullName is provided, update the user
	if userData.FullName != "" {
		user.FullName = userData.FullName
	}

	// If Dob is provided, update the user
	if userData.Dob != (time.Time{}) {
		user.Dob = &userData.Dob
	}

	// If LinkedinURL is provided, update the user
	if userData.LinkedinURL != "" {
		user.LinkedinURL = userData.LinkedinURL
	}

	// Save the user updates to the database
	result = storage.DB.Save(&user)

	// If there is an error, we return a bad request
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	// Return a success response with the updated user
	response.SuccessJSON(c, http.StatusOK, "User updated successfully!", &user)
}

/**
 * Function to delete user by id
 * @api {delete} /users/:id Delete user by id
 * @apiSuccessExample {json} Success-Response:
 * HTTP/1.1 200 OK
 * {
 * 	"success": true,
 * 	"message": "User deleted successfully!",
 * 	"data": null
 * }
 * @apiErrorExample {json} Error-Response:
 * HTTP/1.1 400 Bad Request
 * {
 * 	"success": false,
 * 	"message": "Error deleting user!",
 * 	"data": null
 * }
 */
func (u *UserController) DeleteUserById(c *gin.Context) {
	user := models.User{}
	// Get the user by id
	result := storage.DB.Where("id = ?", c.Param("id")).First(&user)
	// If user not found, return a bad request
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}
	// Soft delete the user
	result = storage.DB.Delete(&user)

	// If there is an error, we return a bad request
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	// Return a success response
	response.SuccessJSON(c, http.StatusOK, "User deleted successfully!", nil)
}
