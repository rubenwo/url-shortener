package server

import (
	"encoding/json"
	"net/http"
)

type JsonError struct {
	StatusCode int
	Msg        string
}

func writeJsonError(w http.ResponseWriter, err error, statusCode int) {
	msg := JsonError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&msg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
