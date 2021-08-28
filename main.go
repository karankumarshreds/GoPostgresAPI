package main

import "os"

const PORT = ":8000"

func main() {
	// creating an app instance 
	app := App{}
	app.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	app.Run(PORT)
}
