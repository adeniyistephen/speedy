package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("finish")
	return nil
}
