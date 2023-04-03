package services

import (
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/HeadHardener/tp_lab/internal/app/repositories"
)

type Authorization interface {
	GenerateToken(workerInput models.LogWorkerInput) (string, error)
	ParseToken(accessToken string) (models.WorkerAttributes, error)
}

type Administration interface {
	CreateWorker(workerInput models.CreateWorkerInput) (int, error)
	GetAll() ([]models.Worker, error)
	GetByID(workerID int) (models.Worker, error)
	UpdateWorker(workerID int, workerInput models.UpdateWorkerInput) error
}

type GSMInterface interface {
	Create(workerID int, docInput models.CreateDocInput) (int, error)
	GetAll() ([]models.Document, error)
	GetByID(docID int) (models.Document, error)
	Update(docID int, docInput models.UpdateDocInput) error
	Delete(docID int) error
}

type Service struct {
	Authorization
	Administration
	GSMInterface
}

func NewService(repos *repositories.Repository) *Service {
	return &Service{
		Authorization:  NewAuthService(repos),
		Administration: NewAdminService(repos),
		GSMInterface:   NewGSMService(repos),
	}
}
