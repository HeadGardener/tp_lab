package repositories

import (
	"fmt"
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/jmoiron/sqlx"
)

type WorkerRepository struct {
	db *sqlx.DB
}

func NewWorkerRepository(db *sqlx.DB) *WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) CreateWorker(worker models.Worker) (int, error) {
	var id int
	query := fmt.Sprintf(`insert into %s (name, surname, fathers_name, phone, role, password_hash)
								values ($1, $2, $3, $4, $5, $6)
								returning id`, workersTable)

	if err := r.db.QueryRow(query,
		worker.Name,
		worker.Surname,
		worker.FathersName,
		worker.Phone,
		worker.Role,
		worker.PasswordHash).
		Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *WorkerRepository) GetWorker(worker *models.Worker) error {
	query := fmt.Sprintf(`select * from %s where name=$1 and surname=$2 and phone=$3 and password_hash=$4`,
		workersTable)

	if err := r.db.Get(worker, query,
		worker.Name,
		worker.Surname,
		worker.Phone,
		worker.PasswordHash); err != nil {
		return err
	}

	return nil
}
