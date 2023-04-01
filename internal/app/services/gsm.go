package services

import (
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/HeadHardener/tp_lab/internal/app/repositories"
)

type GSMService struct {
	repos *repositories.Repository
}

func NewGSMService(repos *repositories.Repository) *GSMService {
	return &GSMService{repos: repos}
}

func (s *GSMService) Create(workerID int, docInput models.CreateDocInput) (int, error) {
	document := models.Document{
		Car:        docInput.Car,
		CarID:      docInput.CarID,
		Waybill:    docInput.Waybill,
		DriverName: docInput.DriverName,
		GasAmount:  docInput.GasAmount,
		GasType:    docInput.GasType,
		IssueDate:  docInput.IssueDate,
	}

	return s.repos.GSMInterface.Create(workerID, document)
}

func (s *GSMService) GetAll() ([]models.Document, error) {
	return s.repos.GSMInterface.GetAll()
}

func (s *GSMService) GetByID(docID int) (models.Document, error) {
	return s.repos.GSMInterface.GetByID(docID)
}
