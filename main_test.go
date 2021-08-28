package main_test 

import (
	"testing"
	"log"
	"os"
	"github.com/karankumarshreds/GoPostgresAPI.git"
)

// application we want to test 
var a main.App

// this will be executed before each tests 
func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
}

// makes sure the table exists by creating one 
func ensureTableExists() {
	tableCreationQuery := `CREATE TABLE IF NOT EXISTS products
	(
		id SERIAL,
		name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id) 
	)
	`
	_, err := a.DB.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}