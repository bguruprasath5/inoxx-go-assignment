package middlewares

import (
	"ionixx/api/models"
	"ionixx/api/response"
	"ionixx/api/storage"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

/**
 * Middleware to check if the user has a valid token to access the route
 */
func AuthMiddleware(c *gin.Context) {
	// Get the token from the header
	idTokenHeader := strings.Split(c.GetHeader("Authorization"), "Bearer ")

	// If the token is not found, return an error
	if len(idTokenHeader) < 2 {
		response.ErrorJSON(c, http.StatusUnauthorized, "Must provide Authorization header with format `Bearer {token}`")
		return
	}
	// If the token is found, get the token
	tokenString := idTokenHeader[1]

	// Parse the token
	token, _ := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Check if the token is valid
	if claims, ok := token.Claims.(*models.JwtCustomClaims); ok && token.Valid {
		// Check if the user is in the database with user id and provided token
		var user models.User
		storage.DB.Model(&models.User{}).Where("auth_token = ? and id = ?", tokenString, claims.UserID).First(&user)
		if user.ID == 0 {
			response.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized: Invalid token")
			return
		}
		// Set user id so that can be accessed from route
		c.Set("userId", user.ID)

		//Proceed to route or next middleware
		c.Next()
	} else {
		//Return error
		response.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized: Invalid token")
		return
	}
}
