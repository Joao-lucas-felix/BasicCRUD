package database

import (
	"database/sql"

	_ "github.com/lib/pq" //Go driver connect
)
//this function try to connect with the Data Base and return the connection
func Connect() (*sql.DB, error) {
	strConnection := "user=golang dbname=devbook password=golang host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", strConnection)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
