package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// ConnectDB returns connection object
func ConnectDB() *sql.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	connStr := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		username,
		password,
		dbPort,
		dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(err.Error())
	}

	createEmployeeTable(db)
	return db
}

// createEmployeeTable creates employees table in DB
func createEmployeeTable(db *sql.DB) {
	const tableCreationQuery = `CREATE TABLE IF NOT EXISTS employees (id SERIAL, name TEXT NOT NULL, position TEXT NOT NULL, salary NUMERIC(10,2) NOT NULL DEFAULT 0.00)`
	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// ClearEmployeeTable clears employees table
func ClearEmployeeTable(db *sql.DB) {
	db.Exec("DELETE FROM employees")
	db.Exec("ALTER SEQUENCE employees_id_seq RESTART WITH 1")
}
