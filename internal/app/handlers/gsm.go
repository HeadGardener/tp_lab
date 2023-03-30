package handlers

import "net/http"

func (h *Handler) createDocument(w http.ResponseWriter, r *http.Request) {
	newResponse(w, http.StatusOK, r.Context().Value(workerCtx))
}

func (h *Handler) getAllDocuments(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) getDocumentByID(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) updateDocument(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) deleteDocument(w http.ResponseWriter, r *http.Request) {
}
