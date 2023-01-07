package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Connect to the database
	var err error
	db, err = sql.Open("mysql", "root:polarbear21*@tcp(localhost:3306)/winter")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set up the Gin router
	r := gin.Default()

	// Handle POST request to submit user's id and password
	r.POST("/submit", func(c *gin.Context) {
		// Bind the JSON payload to a struct
		var payload struct {
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save the user's id and password to the database
		_, err := db.Exec("INSERT INTO clients (name, password) VALUES (?, ?)", payload.Name, payload.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Handle GET request to retrieve user's information
	r.GET("/user/:name", func(c *gin.Context) {
		// Get the user's name from the URL parameters
		name := c.Param("name")

		// Query the database for the user's information
		row := db.QueryRow("SELECT id, password FROM clients WHERE name = ?", name)

		// Scan the result
		var id int
		var password string
		err := row.Scan(&id, &password)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// Return the user's information
		c.JSON(http.StatusOK, gin.H{
			"id":       id,
			"name":     name,
			"password": password,
		})
	})

	// Start the server
	r.Run(":8080")
}
