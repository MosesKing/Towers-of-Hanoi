package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv" // Import strconv package for string conversion utilities.
)

// Move represents a single operation in the Tower of Hanoi puzzle.
type Move struct {
	Disc int    `json:"disc"` // The disc number being moved.
	From string `json:"from"` // The rod from which the disc is moved.
	To   string `json:"to"`   // The rod to which the disc is moved.
}

// solveHanoi recursively calculates the moves required to solve the Tower of Hanoi puzzle.
func solveHanoi(n int, from, to, aux string, moves *[]Move) {
	if n == 1 {
		// Base case: only one disc to move directly from the source to destination.
		*moves = append(*moves, Move{Disc: n, From: from, To: to})
		return
	}
	// Recursive case: Move n-1 discs to the auxiliary rod.
	solveHanoi(n-1, from, aux, to, moves)
	// Move the nth disc to the destination rod.
	*moves = append(*moves, Move{Disc: n, From: from, To: to})
	// Move the n-1 discs from the auxiliary rod to the destination rod.
	solveHanoi(n-1, aux, to, from, moves)
}

// handler handles the HTTP requests for solving the Tower of Hanoi puzzle.
func handler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the number of discs from the query parameter "discs"
	discsStr := r.URL.Query().Get("discs")
	if discsStr == "" {
		// Respond with an error if the "discs" parameter is missing.
		http.Error(w, "Missing discs parameter", http.StatusBadRequest)
		return
	}

	// Convert the discs string to an integer
	discs, err := strconv.Atoi(discsStr)
	if err != nil {
		// Respond with an error if the conversion fails (non-integer value).
		http.Error(w, "Invalid number of discs: must be an integer", http.StatusBadRequest)
		return
	}

	// Ensure the number of discs is a positive integer
	if discs <= 0 {
		// Respond with an error if the number of discs is not positive.
		http.Error(w, "Invalid number of discs: must be a positive integer", http.StatusBadRequest)
		return
	}

	// Compute the moves for the specified number of discs
	var moves []Move
	solveHanoi(discs, "A", "C", "B", &moves)

	// Set the response header to application/json and encode the moves into JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moves)
}

// main initializes the HTTP server.
func main() {
	// Set up HTTP routing and print server start-up message.
	http.HandleFunc("/", handler)
	fmt.Println("Starting server at port 8080...")
	// Start an HTTP server listening on port 8080 and log any errors.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
