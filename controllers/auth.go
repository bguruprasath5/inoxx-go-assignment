package controllers

import (
	"ionixx/constants"
	"ionixx/models"
	"ionixx/storage"
	"ionixx/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthController struct{}

func generateToken(user models.User) (string, error) {
	claims := &models.JwtCustomClaims{
		UserID:   user.ID,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    "ionixx",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(constants.HmacSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthController) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User

	result := storage.DB.Where("user_name = ?", loginRequest.UserName).First(&user)
	if result.Error != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, "Invalid username or password")
		return
	}

	token, err := generateToken(user)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	result = storage.DB.Model(&user).Update("auth_token", token)
	if result.Error != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	utils.SuccessJSON(c, http.StatusOK, "Login successful!", user)
}
