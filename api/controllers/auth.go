package controllers

import (
	"ionixx/api/models"
	"ionixx/api/response"
	"ionixx/api/storage"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

/**
 * Function to generate Jwt token tith custom User claim
 */
func generateToken(user models.User) (string, error) {
	//Creating custom jwt claims from user model
	claims := &models.JwtCustomClaims{
		UserID:   user.ID,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    "ionixx",
			IssuedAt:  time.Now().Unix(),
		},
	}

	//Creating new jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Signing the jwt token with the secret
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/**
 * Function to handle the login request and generate Jwt token
 * @api {post} /auth/login Login
 * @apiSuccessExample {json} Success-Response:
 *   HTTP/1.1 200 OK
 *  {
 *   	"success": true,
 *   	"message": "Login successful!",
 * 		"data": {
 * 			"id": 1,
 * 			"user_name": "test",
 * 			"full_name": "Test User",
 * 			"auth_token: "",
 * 			"dob": "2020-01-01T00:00:00Z",
 * 			"linkedin_url": "",
 * 			"created_at": "2020-01-01T00:00:00Z",
 * 			"updated_at": "2020-01-01T00:00:00Z",
 * 		}
 * 	}
 * @apiErrorExample {json} Error-Response:
 * HTTP/1.1 400 Bad Request
 * {
 * 	"success": false,
 * 	"message": "Error logging in!",
 * 	"data": null
 * }
 */
func (a *AuthController) Login(c *gin.Context) {

	var loginRequest LoginRequest
	//Bind the request body to the login request
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		//If there is an error return bad request
		response.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	//Find the user with the given username
	result := storage.DB.Where("user_name = ?", loginRequest.UserName).First(&user)
	//If the user is not found return bad request
	if result.Error != nil {
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}
	//If the user is found check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		//If the password is incorrect return bad request
		response.ErrorJSON(c, http.StatusBadRequest, "Invalid username or password")
		return
	}
	//If the password is correct generate the jwt token
	token, err := generateToken(user)
	if err != nil {
		//If there is an error return bad request
		response.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}
	//If token is generated successfully update the user with the token
	result = storage.DB.Model(&user).Update("auth_token", token)
	if result.Error != nil {
		//If there is an error return bad request
		response.ErrorJSON(c, http.StatusBadRequest, result.Error.Error())
		return
	}
	//Login successful return the user
	response.SuccessJSON(c, http.StatusOK, "Login successful!", user)
}
