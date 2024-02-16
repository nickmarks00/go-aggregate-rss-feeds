package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Responding with 500 level error: ", message)
	}

	type errResponse struct {
		Error string `json:"error"` // tells payload to marshall into JSON obj with error key
	}

	respondWithJSON(w, code, errResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload) // returns the data as bytes
	if err != nil {
		log.Printf("Failed to marshall JSON response: %v", payload)
		w.WriteHeader(500) // internal server error
		return
	}

	w.Header().Add("Content-Type", "application/json") // add appropriate response header
	w.WriteHeader(code)                                // success
	w.Write(dat)                                       // write response body
}
