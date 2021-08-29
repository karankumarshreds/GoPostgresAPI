package main

import (
	"database/sql"
	"errors"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"string"` 
	Price float64 `json:"price"`
}

func (p *Product) getProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *Product) updateProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *Product) deleteProduct(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (p *Product) createProduct(db *sql.DB) error {
  return errors.New("Not implemented")
}

func getProducts(db *sql.DB, start, count int) ([]Product, error) {
	rows, err := db.Query(
		"SELECT id, name, price FROM products LIMIT $1 OFFSET $2",
		count, start,
	)
	if err != nil {
		return nil, err 
	}
	// Close closes the database and prevents new queries from starting. 
	// Close then waits for all queries that have started processing on 
	// the server to finish.
	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var p Product 
		// can copies the columns from the matched row into the values pointed at by dest
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		} else {
			products = append(products, p)
		}
	}
	return products, nil
}