package repositories

import (
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/jmoiron/sqlx"
)

type WorkerInterface interface {
	CreateWorker(worker models.Worker) (int, error)
	GetWorker(worker *models.Worker) error
}

type Repository struct {
	WorkerInterface
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		WorkerInterface: NewWorkerRepository(db),
	}
}
