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

func (u *UserController) GetAllUsers(c *gin.Context) {
	var list []models.User
	result := storage.DB.Find(&list)

	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	response.SuccessJSON(c, http.StatusOK, "Users fetched successfully!", list)
}

func (u *UserController) GetUserByID(c *gin.Context) {
	var user models.User

	result := storage.DB.Where("id = ?", c.Param("id")).First(&user)

	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

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

func (u *UserController) CreateUser(c *gin.Context) {
	var userData CreateUserRequest

	if err := c.ShouldBindJSON(&userData); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			response.ErrorJSON(c, http.StatusBadRequest, fmt.Sprintf("Validation Error : %s %s", fieldErr.Field(), fieldErr.Tag()))
			return
		}
		return
	}

	user, userErr := userData.toUser()

	if userErr != nil {
		response.ErrorJSON(c, http.StatusBadRequest, userErr.Error())
		return
	}
	result := storage.DB.Create(user)

	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	response.SuccessJSON(c, http.StatusOK, "User created successfully!", nil)

}

type UpdateUserRequest struct {
	FullName    string    `json:"full_name"`
	Dob         time.Time `json:"dob"`
	LinkedinURL string    `json:"linkedin_url"`
}

func (u *UserController) UpdateUserById(c *gin.Context) {
	var userData UpdateUserRequest

	if err := c.BindJSON(&userData); err != nil {
		response.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	user := models.User{}
	result := storage.DB.Where("id = ?", c.Param("id")).First(&user)
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	if userData.FullName != "" {
		user.FullName = userData.FullName
	}
	if userData.Dob != (time.Time{}) {
		user.Dob = &userData.Dob
	}
	if userData.LinkedinURL != "" {
		user.LinkedinURL = userData.LinkedinURL
	}
	result = storage.DB.Save(&user)
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}
	response.SuccessJSON(c, http.StatusOK, "User updated successfully!", nil)
}

func (u *UserController) DeleteUserById(c *gin.Context) {
	user := models.User{}
	result := storage.DB.Where("id = ?", c.Param("id")).First(&user)
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}
	result = storage.DB.Delete(&user)
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}
	response.SuccessJSON(c, http.StatusOK, "User deleted successfully!", nil)
}
