package handlers

import "net/http"

func (h *Handler) isValid(w http.ResponseWriter, r *http.Request) {
	newResponse(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	workerID, err := getWorkerID(r)
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	worker, err := h.service.Administration.GetByID(workerID)
	if err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(w, http.StatusOK, worker)
}
