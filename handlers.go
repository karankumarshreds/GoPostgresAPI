package main 

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"encoding/json"
)

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return 
	}
	p := Product{ID: id}
	err = p.getProduct(a.DB)
	if err != nil {
		switch err {
			case sql.ErrNoRows:
				respondWithError(w, http.StatusNotFound, "Product not found")
			default:
				respondWithError(w, http.StatusBadGateway, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}


// **********************
//    HELPER METHODS 
// **********************

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	var errorMessage  = map[string] string{"error": message}
	respondWithJSON(w, statusCode, errorMessage)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	// encode the response for the client 
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	// WriteHeader sends an HTTP response header with the provided
	// status code.
	w.WriteHeader(statusCode)
	w.Write(response)
}