package middlewares

import (
	"ionixx/api/constants"
	"ionixx/api/models"
	"ionixx/api/response"
	"ionixx/api/storage"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(c *gin.Context) {
	idTokenHeader := strings.Split(c.GetHeader("Authorization"), "Bearer ")

	if len(idTokenHeader) < 2 {
		response.ErrorJSON(c, http.StatusUnauthorized, "Must provide Authorization header with format `Bearer {token}`")
		return
	}
	tokenString := idTokenHeader[1]
	token, _ := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return constants.HmacSecret, nil
	})

	if claims, ok := token.Claims.(*models.JwtCustomClaims); ok && token.Valid {
		var user models.User
		storage.DB.Model(&models.User{}).Where("auth_token = ? and id = ?", tokenString, claims.UserID).First(&user)
		if user.ID == 0 {
			response.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized: Invalid token")
			return
		}
		c.Set("userId", user.ID)
		c.Next()
	} else {
		response.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized: Invalid token")
		return
	}
}
