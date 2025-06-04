package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var shouldFail bool

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Only PATCH is allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Println("Received PATCH:", payload)

	if shouldFail {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payload received",
		"payload": payload,
	})
}

func main() {
	if os.Getenv("MOCK_SHOULD_FAIL") == "true" {
		shouldFail = true
	}

	http.HandleFunc("/api/v1", handler)

	log.Println("Mock API running on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
