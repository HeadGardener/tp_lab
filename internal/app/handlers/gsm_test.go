package handlers

import (
	"bytes"
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
	"time"
)

var toMyTime = func(ts string) models.MyTime {
	t, _ := time.Parse("2006-01-02", ts)
	return models.MyTime(t)
}

func TestHandler_createDocument(t *testing.T) {
	type mockBehavior func(s *mock_services.MockGSMInterface, document models.CreateDocInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputDocument        models.CreateDocInput
		workerAtr            models.WorkerAttributes
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			inputBody: `{"car":"test_car", "car_id":"1111 AA-1", "waybill": 1111, "driver_name":"test_name", 
						"gas_amount": 1, "gas_type":"95", "issue_date":"2023-01-01"}`,
			inputDocument: models.CreateDocInput{
				Car:        "test_car",
				CarID:      "1111 AA-1",
				Waybill:    1111,
				DriverName: "test_name",
				GasAmount:  1,
				GasType:    "95",
				IssueDate:  toMyTime("2023-01-01"),
			},
			workerAtr: models.WorkerAttributes{
				ID:   1,
				Role: "admin",
				Name: "Test",
			},
			mockBehavior: func(s *mock_services.MockGSMInterface, document models.CreateDocInput) {
				s.EXPECT().Create(1, document).Return(1, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: "{\"documentID\":1}\n",
		},
		{
			name: "invalid car_id",
			inputBody: `{"car":"test_car", "car_id":"111 AA-1", "waybill": 1111, "driver_name":"test_name", 
						"gas_amount": 1, "gas_type":"95", "issue_date":"2023-01-01"}`,
			workerAtr: models.WorkerAttributes{
				ID:   1,
				Role: "admin",
				Name: "Test",
			},
			mockBehavior:         func(s *mock_services.MockGSMInterface, document models.CreateDocInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\":\"invalid car_id\"}\n",
		},
		{
			name: "invalid waybill",
			inputBody: `{"car":"test_car", "car_id":"1111 AA-1", "waybill": 111, "driver_name":"test_name", 
						"gas_amount": 1, "gas_type":"95", "issue_date":"2023-01-01"}`,
			workerAtr: models.WorkerAttributes{
				ID:   1,
				Role: "admin",
				Name: "Test",
			},
			mockBehavior:         func(s *mock_services.MockGSMInterface, document models.CreateDocInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\":\"invalid waybill value\"}\n",
		},
		{
			name: "service failure",
			inputBody: `{"car":"test_car", "car_id":"1111 AA-1", "waybill": 1111, "driver_name":"test_name", 
						"gas_amount": 1, "gas_type":"95", "issue_date":"2023-01-01"}`,
			inputDocument: models.CreateDocInput{
				Car:        "test_car",
				CarID:      "1111 AA-1",
				Waybill:    1111,
				DriverName: "test_name",
				GasAmount:  1,
				GasType:    "95",
				IssueDate:  toMyTime("2023-01-01"),
			},
			workerAtr: models.WorkerAttributes{
				ID:   1,
				Role: "admin",
				Name: "Test",
			},
			mockBehavior: func(s *mock_services.MockGSMInterface, document models.CreateDocInput) {
				s.EXPECT().Create(1, document).Return(0, errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "{\"message\":\"service failure\"}\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			gsmService := mock_services.NewMockGSMInterface(c)
			tc.mockBehavior(gsmService, tc.inputDocument)

			service := &services.Service{GSMInterface: gsmService}
			handler := NewHandler(service)

			router := chi.NewRouter()
			router.Post("/api/gsm/", handler.createDocument)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/gsm/", bytes.NewBufferString(tc.inputBody))
			r = r.WithContext(context.WithValue(r.Context(), workerCtx, tc.workerAtr))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getDocumentByID(t *testing.T) {
	type mockBehavior func(s *mock_services.MockGSMInterface, docID any)

	testTable := []struct {
		name                 string
		docID                any
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:  "ok",
			docID: 1,
			mockBehavior: func(s *mock_services.MockGSMInterface, docID any) {
				s.EXPECT().GetByID(docID).Return(models.Document{
					ID:         1,
					Car:        "test_car",
					CarID:      "1111 AA-1",
					Waybill:    1111,
					DriverName: "test_name",
					GasAmount:  1,
					GasType:    "95",
					IssueDate:  toMyTime("2023-01-01"),
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: fmt.Sprintf("{\"ID\":1,\"car\":\"test_car\"," +
				"\"car_id\":\"1111 AA-1\",\"waybill\":1111,\"driver_name\":\"test_name\"," +
				"\"gas_amount\":1,\"gas_type\":\"95\",\"issue_date\":\"2023-01-01\"}\n"),
		},
		{
			name:                 "invalid document_id param",
			docID:                "bad_id",
			mockBehavior:         func(s *mock_services.MockGSMInterface, docID any) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\":\"invalid document_id param\"}\n",
		},
		{
			name:  "service failure",
			docID: 1,
			mockBehavior: func(s *mock_services.MockGSMInterface, docID any) {
				s.EXPECT().GetByID(docID).Return(models.Document{}, errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "{\"message\":\"service failure\"}\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			gsmService := mock_services.NewMockGSMInterface(c)
			tc.mockBehavior(gsmService, tc.docID)

			service := &services.Service{GSMInterface: gsmService}
			handler := NewHandler(service)

			router := chi.NewRouter()
			router.Get("/api/gsm/{document_id}", handler.getDocumentByID)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", fmt.Sprintf("/api/gsm/%v", tc.docID), bytes.NewBufferString(""))
			// r = r.WithContext(context.WithValue(r.Context(), workerCtx, tc.workerAtr))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
