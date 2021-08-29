package main

import (
	"fmt"
	"log"
	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// struct named App to hold our application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// method to initialize the app instance 
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	var err error 
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Cannot connect to the database right now")
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// methods to initialize all the routes 
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")	
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

// method to run the app instance 
func (a *App) Run(addr string) {
	
}
