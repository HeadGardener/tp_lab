package handlers

import (
	"bytes"
	"errors"
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

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuthorization, worker models.LogWorkerInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputWorker          models.LogWorkerInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name":"Vitali", "surname":"Tsaal", "phone":"+375 44 571-19-05", "password":"12345"}`,
			inputWorker: models.LogWorkerInput{
				Name:     "Vitali",
				Surname:  "Tsaal",
				Phone:    "+375 44 571-19-05",
				Password: "12345",
			},
			mockBehavior: func(s *mock_services.MockAuthorization, worker models.LogWorkerInput) {
				s.EXPECT().GenerateToken(worker).Return("token", nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: "{\"token\":\"token\"}\n",
		},
		{
			name:                 "empty name field",
			inputBody:            `{"surname":"Tsaal", "phone":"+375 44 571-19-05", "password":"12345"}`,
			mockBehavior:         func(s *mock_services.MockAuthorization, worker models.LogWorkerInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\":\"name fields can't be empty\"}\n",
		},
		{
			name:                 "invalid phone",
			inputBody:            `{"name":"Vitali", "surname":"Tsaal", "phone":"+375 44571-19-05", "password":"12345"}`,
			mockBehavior:         func(s *mock_services.MockAuthorization, worker models.LogWorkerInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\":\"empty or invalid phone number\"}\n",
		},
		{
			name:      "service failure",
			inputBody: `{"name":"Vitali", "surname":"Tsaal", "phone":"+375 44 571-19-05", "password":"12345"}`,
			inputWorker: models.LogWorkerInput{
				Name:     "Vitali",
				Surname:  "Tsaal",
				Phone:    "+375 44 571-19-05",
				Password: "12345",
			},
			mockBehavior: func(s *mock_services.MockAuthorization, worker models.LogWorkerInput) {
				s.EXPECT().GenerateToken(worker).Return("", errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "{\"message\":\"service failure\"}\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_services.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputWorker)

			service := &services.Service{Authorization: auth}
			handler := NewHandler(service)

			router := chi.NewRouter()
			router.Post("/api/auth/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/auth/sign-in", bytes.NewBufferString(tc.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
