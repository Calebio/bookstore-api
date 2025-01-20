package main

import (
	"bookstore-api/data"
	"bookstore-api/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize database connection
	dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	data.InitDB(dsn)
	defer data.CloseDB()

	// Set up routes
	http.HandleFunc("/books", handlers.GetBooks)          // GET all books
	http.HandleFunc("/book", handlers.GetBook)           // GET book by ID
	http.HandleFunc("/book/create", handlers.CreateBook) // POST create book
	http.HandleFunc("/book/update", handlers.UpdateBook) // PUT update book
	http.HandleFunc("/book/delete", handlers.DeleteBook) // DELETE book

	// Start the server
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
