package model

import (
	"database/sql"
  "fmt"
)

type DB struct {
	*sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "ucactus"
	password = "hnbcactus"
	dbname   = "dbcactus"
)
var db *sql.DB

func InitDB() (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
