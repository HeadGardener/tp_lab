package handlers

import (
	"encoding/json"
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"net/http"
)

func (h *Handler) createDocument(w http.ResponseWriter, r *http.Request) {
	var docInput models.CreateDocInput

	if err := json.NewDecoder(r.Body).Decode(&docInput); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := docInput.Validate(); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	workerID, err := getWorkerID(r)
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	docID, err := h.service.GSMInterface.Create(workerID, docInput)
	if err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(w, http.StatusCreated, map[string]interface{}{
		"documentID": docID,
	})
}

func (h *Handler) getAllDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := h.service.GSMInterface.GetAll()
	if err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(w, http.StatusOK, documents)
}

func (h *Handler) getDocumentByID(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) updateDocument(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) deleteDocument(w http.ResponseWriter, r *http.Request) {
}
