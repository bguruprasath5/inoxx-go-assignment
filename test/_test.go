package test

import (
	"bytes"
	"encoding/json"
	"ionixx/api/controllers"
	"ionixx/api/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupTestServer() *gin.Engine {
	// Switch to test mode
	gin.SetMode(gin.TestMode)

	storage.InitDB()
	// Setup router, just like main function
	r := gin.Default()
	return r
}

func GetUserPayload() ([]byte, error) {
	newUser := controllers.CreateUserRequest{
		UserName:    "tester",
		Password:    "test#123",
		FullName:    "Tester",
		Dob:         time.Now(),
		LinkedinURL: "https://www.linkedin.com/in/tester",
	}
	return json.Marshal(newUser)
}

func GetInvalidUserPayload() ([]byte, error) {
	newUser := controllers.CreateUserRequest{
		UserName:    "tester",
		Password:    "tes",
		FullName:    "Tester",
		LinkedinURL: "https://www.linkedin.com/in/tester",
	}
	return json.Marshal(newUser)
}

func SeedNewUser(r *gin.Engine, t *testing.T) {
	userController := controllers.UserController{}
	r.POST("/", userController.CreateUser)

	data, err := GetUserPayload()

	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(data))

	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder to inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)
}
