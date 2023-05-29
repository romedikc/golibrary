package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 10},
}

// get request
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books) // get data serialized into json
}

// actual detail view where we pass the helper func
func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

// patch request
func checkoutBook(c *gin.Context) {
	// getting id and checking if it exists
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}
	// getting the book and checking if it exists
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	// buying the book and reduce the quantity
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

// patch request
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

// helper func to get id of the book
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

// post request
func postBook(c *gin.Context) {
	var newBook book // book that is to be created

	// error handling
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	// if not error, we append (post) a new book to the list
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", postBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")

}
