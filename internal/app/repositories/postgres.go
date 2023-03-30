package repositories

import (
	"fmt"
	"github.com/HeadHardener/tp_lab/configs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	workersTable     = "workers"
	docsTable        = "documents"
	workersDocsTable = "workers_documents"
)

func NewDB(conf configs.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx",
		fmt.Sprintf("host=%s dbname=%s sslmode=%s", conf.Host, conf.DBName, conf.SSLMode))

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
