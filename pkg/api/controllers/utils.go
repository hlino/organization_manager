package controllers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Errorf("could not serialize json: %v: %v", err, data)
		statusCode = http.StatusInternalServerError
		data = ErrorResponse{err.Error()}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}
