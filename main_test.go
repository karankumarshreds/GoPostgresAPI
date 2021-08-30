package main_test 

import (
	"bytes"
	"testing"
	"encoding/json"
	"strconv"
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
	req, _ := http.NewRequest("GET", "/products?count=5&skip=5", nil)
	// rr := executeRequest(req)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)	
	
	// checking the response status code 
	if http.StatusOK != rr.Code {
		// this method ends the running tests 
		t.Errorf("Expected response %v got %v", http.StatusOK, rr.Code)
	}

	// checking the body response
	body := rr.Body.String();
	if body != "[]" {
		t.Errorf("Expected an empty array. Got %v", body)
	}

}	

func TestGetNonExistentProduct(t *testing.T) {

	clearTable()

	req, _ := http.NewRequest("GET", "/product/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected response %v got %v", http.StatusNotFound, rr.Code)
	}

	var m map[string] string 
	// decode the json to the map
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error key of the response to set to 'Product not found got %v", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {

	clearTable()
	var jsonStr = []byte(`{"name": "test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected response code %v got %v", http.StatusCreated, rr.Code)
	}

	var m map[string] interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}
	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}
	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
			t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {

	clearTable()

	addProducts(1)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected response %v got %v", http.StatusOK, rr.Code)
	}

	var m map[string] string
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["name"] != "test product1" {
		t.Errorf("Expected product name to be 'test product' got %v", m["name"])
	}

}

func TestUpdateProduct(t *testing.T) {

	clearTable()
	addProducts(1)

	jsonStr := []byte(`{"price":1000}`)
	req, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("Expected response %v got %v", http.StatusOK, rr.Code)
	}

	var m map[string] interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)

	if m["price"] != float64(1000) {
		t.Errorf("Expected product price to be 1000 got %v", m["price"])
	}
	if m["id"] != float64(1) {
		t.Errorf("Expected product ID to be 1 got %v", m["id"])
	}

}

func TestDeleteProduct(t *testing.T) {

	clearTable()
	addProducts(1)
	req, _ := http.NewRequest("DELETE", "/product/1", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected response %v got %v", http.StatusOK, rr.Code)
	}
	// making a get request now to check if the product has been deleted or not 
	req, _ = http.NewRequest("GET", "/product/1", nil)
	rr = httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected response %v got %v", http.StatusNotFound, rr.Code)
	}

}

// *************************** // 
//      HELPER FUNCTIONS 
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
// rw here is the "rr" (response recorder) to record the response 
// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	a.Router.ServeHTTP(rr, req)	
// 	return rr
// }

func addProducts(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "test product" + strconv.Itoa(count), 100)
	}
}