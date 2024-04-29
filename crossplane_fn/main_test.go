package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandler tests the HTTP handler for different scenarios.
func TestHandler(t *testing.T) {
	// Test cases for various scenarios including boundary values and error cases.
	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedMoves  int
	}{
		{"Valid request", "discs=3", http.StatusOK, 7},
		{"Missing discs parameter", "", http.StatusBadRequest, 0},
		{"Non-integer discs", "discs=abc", http.StatusBadRequest, 0},
		{"Zero discs", "discs=0", http.StatusBadRequest, 0},
		{"Negative discs", "discs=-1", http.StatusBadRequest, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", "/?"+test.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Record the response using a response recorder.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handler)

			// Dispatch the request to the handler.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedStatus)
			}

			// If the status is OK, decode the response and check the number of moves.
			if test.expectedStatus == http.StatusOK {
				var moves []Move
				if err := json.NewDecoder(rr.Body).Decode(&moves); err != nil {
					t.Fatal("Could not decode the response:", err)
				}
				if len(moves) != test.expectedMoves {
					t.Errorf("handler returned unexpected number of moves: got %v want %v",
						len(moves), test.expectedMoves)
				}
			}
		})
	}
}
