package test

import (
	"bytes"
	"fmt"
	"ionixx/api/controllers"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
 * Function to test the create user
 * case: success
 */
func TestCreateUserSuccess(t *testing.T) {
	r := SetupTestServer()
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
	fmt.Println(w.Body)

	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}

/**
 * Function to test the create user
 * case: failed
 */
func TestCreateUserFailed(t *testing.T) {

	r := SetupTestServer()
	userController := controllers.UserController{}
	r.POST("/", userController.CreateUser)

	data, err := GetInvalidUserPayload()

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

/**
 * Function to test the get all users
 * case: success
 */
func TestGetAllUsers(t *testing.T) {
	r := SetupTestServer()
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

/**
 * Function to test get user by id endpoint
 * case: sucess
 */
func TestGetUserByIDSuccess(t *testing.T) {

	r := SetupTestServer()
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

/**
 * Function to test get user by id endpoint
 * case: fail
 */
func TestGetUserByIDFailed(t *testing.T) {

	r := SetupTestServer()
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
