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
 * Function to test the login endpoint
 * case: login with valid credentials
 */
func TestLoginSuccess(t *testing.T) {

	// Setup the test server
	r := SetupTestServer()

	// Create a new controller
	authController := controllers.AuthController{}

	// Setup the endpoint
	r.POST("/login", authController.Login)

	// Seed a new user
	SeedNewUser(r, t)

	// Get the payload
	data, err := GetUserPayload()
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	// Create a new request
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))

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
 * Function to test the login endpoint
 * case: login with invalid credentials
 */
func TestLoginFailed(t *testing.T) {

	// Setup the test server
	r := SetupTestServer()
	// Create a new controller
	authController := controllers.AuthController{}
	// Setup the endpoint
	r.POST("/login", authController.Login)

	// Seed a new user
	SeedNewUser(r, t)

	// Get the payload
	data, err := GetInvalidUserPayload()

	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a new request
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))

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
