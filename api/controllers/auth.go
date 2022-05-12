package controllers

import (
	"ionixx/api/constants"
	"ionixx/api/models"
	"ionixx/api/response"
	"ionixx/api/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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
		response.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User

	result := storage.DB.Where("user_name = ?", loginRequest.UserName).First(&user)
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		response.ErrorJSON(c, http.StatusBadRequest, "Invalid username or password")
		return
	}

	token, err := generateToken(user)
	if err != nil {
		response.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	result = storage.DB.Model(&user).Update("auth_token", token)
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}

	response.SuccessJSON(c, http.StatusOK, "Login successful!", user)
}
