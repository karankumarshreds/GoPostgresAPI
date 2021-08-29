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
	return db.QueryRow(
		"SELECT name, price FROM products WHERE id =$1", 
		p.ID,
	).Scan(&p.Name, &p.Price) // scan method only works on method that return row(s)
}

func (p *Product) updateProduct(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE products SET name=$1, price=$2 WHERE id=$3",
		p.Name, 
		p.Price,
		p.ID,
	)
	return err
}

func (p *Product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec(
		"DELETE FROM products WHERE id=$1",
		p.ID,
	)
	return err
}

func (p *Product) createProduct(db *sql.DB) error {
  return db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING name, price, id",
		p.Name, 
		p.Price,
	).Scan(&p.Name, &p.Price, &p.ID)
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
		} 
		products = append(products, p)
		
	}
	return products, nil
}