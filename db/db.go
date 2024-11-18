package db

import (
	"database/sql"
	"log"

	"github.com/MatthewAraujo/vacation-backend/pkg/assert"
	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	assert.NotNil(cfg, "Mysql Config cannot be nil")
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
