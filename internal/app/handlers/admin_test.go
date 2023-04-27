package handlers

import (
	"bytes"
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

func TestHandler_createWorker(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAdministration, worker models.CreateWorkerInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputWorker          models.CreateWorkerInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			inputBody: `{"name":"Test", "surname":"Tested", "fathers_name":"Tester", "phone":"+111 11 111-11-11", "role":"worker",
						"password":"12345"}`,
			inputWorker: models.CreateWorkerInput{
				Name:        "Test",
				Surname:     "Tested",
				FathersName: "Tester",
				Phone:       "+111 11 111-11-11",
				Role:        "worker",
				Password:    "12345",
			},
			mockBehavior: func(s *mock_services.MockAdministration, worker models.CreateWorkerInput) {
				s.EXPECT().CreateWorker(worker).Return(1, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: "{\"id\":1}\n",
		},
		{
			name: "invalid phone",
			inputBody: `{"name":"Test", "surname":"Tested", "fathers_name":"Tester", "phone":"+111 a 111-11-11", "role":"worker",
						"password":"12345"}`,
			inputWorker: models.CreateWorkerInput{
				Name:        "Test",
				Surname:     "Tested",
				FathersName: "Tester",
				Phone:       "+111 a 111-11-11",
				Role:        "worker",
				Password:    "12345",
			},
			mockBehavior:         func(s *mock_services.MockAdministration, worker models.CreateWorkerInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\":\"empty or invalid phone number\"}\n",
		},
		{
			name: "invalid role",
			inputBody: `{"name":"Test", "surname":"Tested", "fathers_name":"Tester", "phone":"+111 11 111-11-11", "role":"boss",
						"password":"12345"}`,
			inputWorker: models.CreateWorkerInput{
				Name:        "Test",
				Surname:     "Tested",
				FathersName: "Tester",
				Phone:       "+111 11 111-11-11",
				Role:        "boss",
				Password:    "12345",
			},
			mockBehavior:         func(s *mock_services.MockAdministration, worker models.CreateWorkerInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\":\"invalid role\"}\n",
		},
		{
			name: "service failure",
			inputBody: `{"name":"Test", "surname":"Tested", "fathers_name":"Tester", "phone":"+111 11 111-11-11", "role":"worker",
						"password":"12345"}`,
			inputWorker: models.CreateWorkerInput{
				Name:        "Test",
				Surname:     "Tested",
				FathersName: "Tester",
				Phone:       "+111 11 111-11-11",
				Role:        "worker",
				Password:    "12345",
			},
			mockBehavior: func(s *mock_services.MockAdministration, worker models.CreateWorkerInput) {
				s.EXPECT().CreateWorker(worker).Return(0, errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "{\"message\":\"service failure\"}\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			adm := mock_services.NewMockAdministration(c)
			tc.mockBehavior(adm, tc.inputWorker)

			service := &services.Service{Administration: adm}
			handler := NewHandler(service)

			router := chi.NewRouter()
			router.Post("/api/admin/worker/sign-up", handler.createWorker)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/admin/worker/sign-up", bytes.NewBufferString(tc.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getWorkerByID(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAdministration, workerID any)

	testTable := []struct {
		name                 string
		workerID             any
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "ok",
			workerID: 1,
			mockBehavior: func(s *mock_services.MockAdministration, workerID any) {
				s.EXPECT().GetByID(workerID).Return(models.Worker{
					ID:           1,
					Name:         "Test",
					Surname:      "Tested",
					FathersName:  "Tester",
					Phone:        "+111 11 111-11-11",
					Role:         "worker",
					PasswordHash: "hash",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: fmt.Sprintf("{\"id\":1,\"name\":\"Test\"," +
				"\"surname\":\"Tested\",\"fathers_name\":\"Tester\",\"phone\":\"+111 11 111-11-11\"," +
				"\"role\":\"worker\",\"password_hash\":\"hash\"}\n"),
		},
		{
			name:                 "invalid worker_id",
			workerID:             "bad_id",
			mockBehavior:         func(s *mock_services.MockAdministration, workerID any) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: fmt.Sprintf("{\"message\":\"invalid worker_id param\"}\n"),
		},
		{
			name:     "service failure",
			workerID: 1,
			mockBehavior: func(s *mock_services.MockAdministration, workerID any) {
				s.EXPECT().GetByID(workerID).Return(models.Worker{}, errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: fmt.Sprintf("{\"message\":\"service failure\"}\n"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			adm := mock_services.NewMockAdministration(c)
			tc.mockBehavior(adm, tc.workerID)

			service := &services.Service{Administration: adm}
			handler := NewHandler(service)

			router := chi.NewRouter()
			router.Get("/api/admin/worker/get/{worker_id}", handler.getWorkerByID)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", fmt.Sprintf("/api/admin/worker/get/%v", tc.workerID), bytes.NewBufferString(""))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
