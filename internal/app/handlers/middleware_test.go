package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/HeadHardener/tp_lab/internal/app/services"
	mock_services "github.com/HeadHardener/tp_lab/internal/app/services/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_identifyUser(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_services.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(models.WorkerAttributes{
					ID:   1,
					Role: "worker",
					Name: "Test",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "no header",
			headerName:           "",
			mockBehavior:         func(s *mock_services.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"message\":\"empty auth header\"}\n",
		},
		{
			name:                 "invalid bearer",
			headerName:           "Authorization",
			headerValue:          "Barer token",
			mockBehavior:         func(s *mock_services.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"message\":\"invalid auth header\"}\n",
		},
		{
			name:                 "no token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(s *mock_services.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"message\":\"jwt token is empty\"}\n",
		},
		{
			name:        "service failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_services.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(models.WorkerAttributes{}, errors.New("failed of token parsing"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: "{\"message\":\"failed of token parsing\"}\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_services.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.token)

			service := &services.Service{Authorization: auth}
			handler := NewHandler(service)

			router := chi.NewRouter()
			router.Use(handler.identifyUser)
			router.Post("/protected", func(w http.ResponseWriter, r *http.Request) {
				workerCtxValue := r.Context().Value(workerCtx)
				workerAttributes, _ := workerCtxValue.(models.WorkerAttributes)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%d", workerAttributes.ID)))
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/protected", nil)
			r.Header.Set(tc.headerName, tc.headerValue)

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_checkRole(t *testing.T) {

	testTable := []struct {
		name                 string
		workerAtr            any
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			workerAtr: models.WorkerAttributes{
				ID:   1,
				Role: "admin",
				Name: "Test",
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "1",
		},
		{
			name:                 "invalid workerAttributes",
			workerAtr:            "invalid Atr",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\":\"workerCtx value is not of type WorkerAttributes\"}\n",
		},
		{
			name: "invalid role",
			workerAtr: models.WorkerAttributes{
				ID:   1,
				Role: "worker",
				Name: "Test",
			},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: "{\"message\":\"you don't have enough rules\"}\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_services.NewMockAuthorization(c)

			service := &services.Service{Authorization: auth}
			handler := NewHandler(service)

			router := chi.NewRouter()
			router.Use(handler.checkRole)
			router.Post("/protected", func(w http.ResponseWriter, r *http.Request) {
				workerCtxValue := r.Context().Value(workerCtx)
				workerAttributes, _ := workerCtxValue.(models.WorkerAttributes)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%d", workerAttributes.ID)))
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/protected", nil)
			ctx := context.WithValue(r.Context(), workerCtx, tc.workerAtr)

			router.ServeHTTP(w, r.WithContext(ctx))

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
