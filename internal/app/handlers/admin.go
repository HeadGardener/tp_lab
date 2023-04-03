package handlers

import (
	"encoding/json"
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) createWorker(w http.ResponseWriter, r *http.Request) {
	var workerInput models.CreateWorkerInput

	if err := json.NewDecoder(r.Body).Decode(&workerInput); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
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
	workers, err := h.service.Administration.GetAll()
	if err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(w, http.StatusOK, workers)
}

func (h *Handler) getWorkerByID(w http.ResponseWriter, r *http.Request) {
	workerID, err := strconv.Atoi(chi.URLParam(r, "worker_id"))
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "invalid worker_id param")
		return
	}

	worker, err := h.service.Administration.GetByID(workerID)
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newResponse(w, http.StatusOK, worker)
}

func (h *Handler) updateWorker(w http.ResponseWriter, r *http.Request) {
	var workerInput models.UpdateWorkerInput

	if err := json.NewDecoder(r.Body).Decode(&workerInput); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	workerID, err := strconv.Atoi(chi.URLParam(r, "worker_id"))
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "invalid worker_id param")
		return
	}

	if err := h.service.Administration.UpdateWorker(workerID, workerInput); err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(w, http.StatusOK, map[string]interface{}{
		"status": "updated",
	})
}

// func (h *Handler) deleteWorker(w http.ResponseWriter, r *http.Request) {}
