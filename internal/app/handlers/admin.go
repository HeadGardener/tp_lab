package handlers

import (
	"encoding/json"
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"net/http"
)

func (h *Handler) createWorker(w http.ResponseWriter, r *http.Request) {
	var workerInput models.CreateWorkerInput

	if err := json.NewDecoder(r.Body).Decode(&workerInput); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "invalid data to decode worker")
		return
	}

	if err := workerInput.Validate(); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Administration.CreateWorker(workerInput)
	if err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(w, http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllWorkers(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getWorkerByID(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) updateWorker(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteWorker(w http.ResponseWriter, r *http.Request) {
}
