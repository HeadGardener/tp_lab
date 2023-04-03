package services

import (
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/HeadHardener/tp_lab/internal/app/repositories"
)

type AdminService struct {
	repos *repositories.Repository
}

func NewAdminService(repos *repositories.Repository) *AdminService {
	return &AdminService{repos: repos}
}

func (s *AdminService) CreateWorker(workerInput models.CreateWorkerInput) (int, error) {
	worker := models.Worker{
		Name:         workerInput.Name,
		Surname:      workerInput.Surname,
		FathersName:  workerInput.FathersName,
		Phone:        workerInput.Phone,
		Role:         workerInput.Role,
		PasswordHash: getPasswordHash(workerInput.Password),
	}

	return s.repos.WorkerInterface.CreateWorker(worker)
}

func (s *AdminService) GetAll() ([]models.Worker, error) {
	return s.repos.WorkerInterface.GetAll()
}

func (s *AdminService) GetByID(workerID int) (models.Worker, error) {
	return s.repos.WorkerInterface.GetByID(workerID)
}

func (s *AdminService) UpdateWorker(workerID int, workerInput models.UpdateWorkerInput) error {
	_, err := s.repos.WorkerInterface.GetByID(workerID)
	if err != nil {
		return err
	}

	worker := workerInput.ToWorker()
	worker.ID = workerID

	return s.repos.WorkerInterface.Update(worker)
}
