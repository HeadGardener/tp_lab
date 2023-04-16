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

func (h *WebSocketHandler) newErrResponse(w http.ResponseWriter, code int, errorMsg string) {
	h.errLogger.Error(errorMsg)
	newResponse(w, code, response{
		Message: errorMsg,
	})
}

func newResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
