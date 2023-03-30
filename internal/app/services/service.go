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
}

type Service struct {
	Authorization
	Administration
}

func NewService(repos *repositories.Repository) *Service {
	return &Service{
		Authorization:  NewAuthService(repos),
		Administration: NewAdminService(repos),
	}
}
