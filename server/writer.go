package server

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, statusCode int, jsonData interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if jsonData != nil {
		return json.NewEncoder(w).Encode(jsonData)
	}
	return nil
}
