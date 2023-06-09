package repositories

import (
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/jmoiron/sqlx"
)

type WorkerInterface interface {
	CreateWorker(worker models.Worker) (int, error)
	GetWorker(worker *models.Worker) error
	GetAll() ([]models.Worker, error)
	GetByID(workerID int) (models.Worker, error)
	Update(worker models.Worker) error
}

type GSMInterface interface {
	Create(workerID int, document models.Document) (int, error)
	GetAll() ([]models.Document, error)
	GetByID(docID int) (models.Document, error)
	GetAllWithWorkerID(workerID int) ([]models.Document, error)
	Update(document models.Document) error
	Delete(docID int) error
}

type Repository struct {
	WorkerInterface
	GSMInterface
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		WorkerInterface: NewWorkerRepository(db),
		GSMInterface:    NewGSMRepository(db),
	}
}
