package handlers

import (
	"bookstore-api/data"
	"bookstore-api/models"
	"encoding/json"
	"net/http"
	"strconv"
)

// GetBooks handles GET requests to retrieve all books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := data.DB.Query("SELECT id, title, author, price FROM books")
	if err != nil {
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
			http.Error(w, "Error scanning books", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// GetBook handles GET requests to retrieve a book by ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	err = data.DB.QueryRow("SELECT id, title, author, price FROM books WHERE id=$1", id).Scan(&book.ID, &book.Title, &book.Author, &book.Price)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// CreateBook handles POST requests to create a new book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := data.DB.QueryRow("INSERT INTO books (title, author, price) VALUES ($1, $2, $3) RETURNING id", book.Title, book.Author, book.Price).Scan(&book.ID)
	if err != nil {
		http.Error(w, "Error creating book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// UpdateBook handles PUT requests to update an existing book
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = data.DB.Exec("UPDATE books SET title=$1, author=$2, price=$3 WHERE id=$4", book.Title, book.Author, book.Price, id)
	if err != nil {
		http.Error(w, "Error updating book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteBook handles DELETE requests to remove a book
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	_, err = data.DB.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}