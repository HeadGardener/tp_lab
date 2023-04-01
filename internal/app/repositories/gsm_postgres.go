package repositories

import (
	"fmt"
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/jmoiron/sqlx"
)

type GSMRepository struct {
	db *sqlx.DB
}

func NewGSMRepository(db *sqlx.DB) *GSMRepository {
	return &GSMRepository{db: db}
}

func (r *GSMRepository) Create(workerID int, document models.Document) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var docID int
	createDocQuery := fmt.Sprintf(`insert into %s 
    									(car, car_id, waybill, driver_name, gas_amount, gas_type, issue_date)
    									values ($1,$2,$3,$4,$5,$6,$7)
    									returning id`, docsTable)

	if err := tx.QueryRow(createDocQuery,
		document.Car,
		document.CarID,
		document.Waybill,
		document.DriverName,
		document.GasAmount,
		document.GasType,
		document.IssueDate).
		Scan(&docID); err != nil {
		tx.Rollback()
		return 0, err
	}

	workersDocsQuery := fmt.Sprintf("insert into %s (worker_id, document_id) values ($1, $2)",
		workersDocsTable)
	if _, err := tx.Exec(workersDocsQuery, workerID, docID); err != nil {
		tx.Rollback()
		return 0, err
	}

	return docID, tx.Commit()
}

func (r *GSMRepository) GetAll() ([]models.Document, error) {
	var documents []models.Document

	query := fmt.Sprintf("select * from %s", docsTable)

	if err := r.db.Select(&documents, query); err != nil {
		return nil, err
	}

	return documents, nil
}
