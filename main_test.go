package main_test 

import (
	"testing"
	"log"
	"fmt"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/karankumarshreds/GoPostgresAPI.git"
)

// application we want to test 
var a main.App

// this will be executed before each tests 
func TestMain(m *testing.M) {
	fmt.Println("Test main function")
	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ensureTableExists()
	
	code := m.Run()

	os.Exit(code)
}



// T is a type passed to Test functions to manage test state and support formatted test logs.
func TestEmptyTable (t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/products", nil)
	// response := executeRequest(req)
	// rw here is the "rr" (response recorder) to record the response 
	rw := httptest.NewRecorder()
	a.Router.ServeHTTP(rw, req)	
	// checking the response status code 
	if http.StatusOK != rw.Code {
		// this method ends the running tests 
		t.Errorf("Expected response %v got %v", http.StatusOK, rw.Code)
	}
	// checking the body response
	body := rw.Body.String();
	if body != "[]" {
		t.Errorf("Expected an empty array. Got %v", body)
	}

}	

// *************************** // 
// HELPER FUNCTIONS 
// *************************** // 

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

// The httptest.ResponseRecorder is an implementation of http.ResponseWriter
// As "rw http.ResponseWriter" assemnles the HTTP server's response by writing to it 
// which we send back the to the HTTP CLient. Similarly, httptest.ResponseRecorder 
// is used to record the response that the handler will write for the client's request 
// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	httptest.NewRecoer
// }