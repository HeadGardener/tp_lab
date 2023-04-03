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

func (r *GSMRepository) GetByID(docID int) (models.Document, error) {
	var document models.Document

	query := fmt.Sprintf("select * from %s where id=$1", docsTable)

	if err := r.db.Get(&document, query, docID); err != nil {
		return models.Document{}, err
	}

	return document, nil
}

func (r *GSMRepository) GetAllWithWorkerID(workerID int) ([]models.Document, error) {
	var documents []models.Document

	query := fmt.Sprintf("select d.* from %s d inner join %s wd on d.id=wd.document_id where wd.worker_id=$1",
		docsTable, workersDocsTable)

	if err := r.db.Select(&documents, query, workerID); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *GSMRepository) Update(document models.Document) error {
	query := fmt.Sprintf(`update %s 
						set car=$1, car_id=$2, waybill=$3, driver_name=$4, gas_amount=$5, gas_type=$6, issue_date=$7
						where id=$8`, docsTable)

	if _, err := r.db.Exec(query,
		document.Car,
		document.CarID,
		document.Waybill,
		document.DriverName,
		document.GasAmount,
		document.GasType,
		document.IssueDate,
		document.ID); err != nil {
		return err
	}

	return nil
}

func (r *GSMRepository) Delete(docID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	deleteFromDocsQuery := fmt.Sprintf("delete from %s where id=$1", docsTable)
	if _, err := tx.Exec(deleteFromDocsQuery, docID); err != nil {
		tx.Rollback()
		return err
	}

	deleteFromWorkersDocs := fmt.Sprintf("delete from %s where document_id=$1", workersDocsTable)
	if _, err := tx.Exec(deleteFromWorkersDocs, docID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
