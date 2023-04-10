package handlers

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func (h *Handler) newErrResponse(w http.ResponseWriter, code int, errorMsg string) {
	h.errLogger.Error(errorMsg)
	newResponse(w, code, response{
		Message: errorMsg,
	})
}

func newResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "authorization")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
