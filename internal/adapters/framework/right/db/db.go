package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type DbAdapter struct {
	db *sql.DB
}

func NewDbAdapter(cfg mysql.Config) (*DbAdapter, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("db connection failure..: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("db ping failure :%v", err)
	}
	return &DbAdapter{db: db}, nil
}

func (da DbAdapter) CloseDbConnection() {
	err := da.db.Close()

	if err != nil {
		log.Fatalf("db close ffailure: %v", err)
	}
}
