package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Close port using `kill -9 $(lsof -t -i:8080)`

// Book Model
type book struct {
	//* Make sure that the name of the field starts with a capital letter.
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// Books Data
var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

// GET /books | curl -X GET 'localhost:8080/books/1'
func getBooks(c *gin.Context) {
	// Return 200
	c.IndentedJSON(http.StatusOK, books)
}

// POST /books | curl -X POST  "localhost:8080/books" --include --header "Content-Type: application/json" -d @body.json
func createBook(c *gin.Context) {
	// `gin.Context` is all the request information (headers, query parametets)
	var newBook book

	// Convers json to object
	if err := c.BindJSON(&newBook); err != nil {
		// Return 400
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request."})
		return
	}

	books = append(books, newBook)

	// Return 201
	c.IndentedJSON(http.StatusCreated, newBook)
}

// GET book/:id | curl -X GET 'localhost:8080/books/1'
func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		// Return 404
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	// Return 200
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not fund")
}

// PATCH /checkout | curl -X PATCH 'localhost:8080/checkout?id=1'
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		// Return 400
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		// Return 404
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		// Return 400
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1
	// Return 200
	c.IndentedJSON(http.StatusOK, book)
}

// PATCH /return | curl -X PATCH 'localhost:8080/return?id=1'
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		// Return 400
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		// Return 404
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1

	// Return 200
	c.IndentedJSON(http.StatusOK, book)

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
