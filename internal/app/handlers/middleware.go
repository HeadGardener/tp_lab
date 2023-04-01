package handlers

import (
	"context"
	"errors"
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"net/http"
	"strings"
)

const (
	workerCtx = "workerAtr"
)

func (h *Handler) identifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			h.newErrResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			h.newErrResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		if headerParts[0] != "Bearer" {
			h.newErrResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		if len(headerParts[1]) == 0 {
			h.newErrResponse(w, http.StatusUnauthorized, "jwt token is empty")
			return
		}

		workerAttributes, err := h.service.Authorization.ParseToken(headerParts[1])
		if err != nil {
			h.newErrResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), workerCtx, workerAttributes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) checkRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workerCtxValue := r.Context().Value(workerCtx)
		workerAttributes, ok := workerCtxValue.(models.WorkerAttributes)
		if !ok {
			h.newErrResponse(w, http.StatusBadRequest, "workerCtx value is not of type WorkerAttributes")
			return
		}

		if workerAttributes.Role != "admin" {
			h.newErrResponse(w, http.StatusForbidden, "you don't have enough rules")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getWorkerID(r *http.Request) (int, error) {
	workerCtxValue := r.Context().Value(workerCtx)
	workerAttributes, ok := workerCtxValue.(models.WorkerAttributes)
	if !ok {
		return 0, errors.New("workerCtx value is not of type WorkerAttributes")
	}

	return workerAttributes.ID, nil
}
