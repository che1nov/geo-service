package controller

import (
	"encoding/json"
	"net/http"
)

type DefaultResponder struct{}

func (dr DefaultResponder) RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (dr DefaultResponder) RespondError(w http.ResponseWriter, status int, message string) {
	dr.RespondJSON(w, status, map[string]string{"error": message})
}
