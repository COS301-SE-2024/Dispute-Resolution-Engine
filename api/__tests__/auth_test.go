package api

import (
    "api/handlers"
    "api/db" // Import the package that contains the db symbo
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

// MockResponseWriter is used to capture the HTTP response
type MockResponseWriter struct {
    httptest.ResponseRecorder
}

func TestCreateUser(t *testing.T) {
	DB := db.Init()
	// Example of a valid request body
    validRequestBody := `{"first_name":"John","surname":"Doe","birthdate":"1990-01-01","nationality":"Country","email":"john.doe@example.com"}`
    // Example of an invalid request body (incomplete JSON)
    invalidRequestBody := `{"first_name":"John"`

    tests := []struct {
        name           string
        body           string
        expectedStatus int
    }{
        {"Valid Request", validRequestBody, http.StatusOK},
        {"Invalid Request", invalidRequestBody, http.StatusBadRequest},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            request, err := http.NewRequest("POST", "/createuser", strings.NewReader(tt.body))
            if err != nil {
                t.Fatal(err)
            }

            // Use httptest.ResponseRecorder to capture the response
            responseRecorder := httptest.NewRecorder()
            handler := handlers.New(DB)

            // Serve the HTTP request to our handler
            handler.CreateUser(responseRecorder, request)

            // Check the status code is what we expect
            if status := responseRecorder.Code; status != tt.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.expectedStatus)
            }
        })
    }
}