package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var books = []Book{
    {ID: "1", Title: "Go Basics", Author: "John Doe"},
    {ID: "2", Title: "Mastering Go", Author: "Jane Smith"},
}

func getBooks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/books/")
    for _, book := range books {
        if book.ID == id {
            json.NewEncoder(w).Encode(book)
            return
        }
    }
    http.NotFound(w, r)
}

func createBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    json.NewDecoder(r.Body).Decode(&book)
    books = append(books, book)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/books/")
    for i, book := range books {
        if book.ID == id {
            json.NewDecoder(r.Body).Decode(&books[i])
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(books[i])
            return
        }
    }
    http.NotFound(w, r)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/books/")
    for i, book := range books {
        if book.ID == id {
            books = append(books[:i], books[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.NotFound(w, r)
}

func main() {
    http.HandleFunc("/books", getBooks)
    http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getBook(w, r)
        case http.MethodPost:
            createBook(w, r)
        case http.MethodPut:
            updateBook(w, r)
        case http.MethodDelete:
            deleteBook(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    fmt.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}