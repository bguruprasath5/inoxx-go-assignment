package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ionixx/api/controllers"
	"ionixx/api/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupTestServer() *gin.Engine {
	// Switch to test mode
	gin.SetMode(gin.TestMode)

	storage.InitDB()
	// Setup router, just like main function
	r := gin.Default()
	return r
}

func getUserPayload() ([]byte, error) {
	newUser := controllers.CreateUserRequest{
		UserName:    "tester",
		Password:    "test#123",
		FullName:    "Tester",
		Dob:         time.Now(),
		LinkedinURL: "https://www.linkedin.com/in/tester",
	}
	return json.Marshal(newUser)
}
func getInvalidUserPayload() ([]byte, error) {
	newUser := controllers.CreateUserRequest{
		UserName:    "tester",
		Password:    "tes",
		FullName:    "Tester",
		LinkedinURL: "https://www.linkedin.com/in/tester",
	}
	return json.Marshal(newUser)
}
func TestCreateUserSuccess(t *testing.T) {

	r := setupTestServer()
	userController := controllers.UserController{}
	r.POST("/", userController.CreateUser)
	data, err := getUserPayload()
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
	fmt.Println(w.Body)

	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}

func TestCreateUserFailed(t *testing.T) {

	r := setupTestServer()
	userController := controllers.UserController{}
	r.POST("/", userController.CreateUser)

	data, err := getInvalidUserPayload()

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
	fmt.Println(w.Body)

	if w.Code == http.StatusBadRequest {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusBadRequest, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusBadRequest, w.Code)
	}
}

func TestGetAllUsers(t *testing.T) {
	r := setupTestServer()
	userController := controllers.UserController{}
	r.GET("/", userController.GetAllUsers)
	SeedNewUser(r, t)

	req, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)
	fmt.Println(w.Body)

	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}

func TestGetUserByIDSuccess(t *testing.T) {

	r := setupTestServer()
	userController := controllers.UserController{}
	r.GET("/:id", userController.GetUserByID)
	SeedNewUser(r, t)

	req, err := http.NewRequest(http.MethodGet, "/1", nil)

	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder to inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)
	fmt.Println(w.Body)

	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}
func TestGetUserByIDFailed(t *testing.T) {

	r := setupTestServer()
	userController := controllers.UserController{}
	r.GET("/:id", userController.GetUserByID)

	req, err := http.NewRequest(http.MethodGet, "/-1", nil)

	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder to inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)
	fmt.Println(w.Body)

	if w.Code == http.StatusBadRequest {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusBadRequest, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusBadRequest, w.Code)
	}
}

func SeedNewUser(r *gin.Engine, t *testing.T) {
	userController := controllers.UserController{}
	r.POST("/", userController.CreateUser)

	data, err := getUserPayload()

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
