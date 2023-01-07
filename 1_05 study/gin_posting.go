package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

// Post represents a single post on the post board.
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Like    int    `json:"like"`
}

// Posts is a slice of Post that represents the post board.
var Posts []Post

func main() {
	router := gin.Default()

	// Show the post board
	router.GET("/posts", func(c *gin.Context) {
		c.JSON(200, Posts)
	})

	// Create a new post
	router.POST("/posts", func(c *gin.Context) {
		var post Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db, err := sql.Open("mysql", "root:polarbear21*@tcp(localhost:3306)/winter")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		defer db.Close()
		// Insert the new post into the database
		result, err := db.Exec("INSERT INTO liketable (Title, Content) VALUES (?, ?)", post.Title, post.Content)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		// Get the ID of the new post
		id, err := result.LastInsertId()

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		post.ID = int(id)
		c.JSON(201, post)
		// Set the post ID
		if len(Posts) > 0 {
			post.ID = Posts[len(Posts)-1].ID + 1
		} else {
			post.ID = 1
		}
		Posts = append(Posts, post)
	})

	// Inquire about a specific post by ID
	router.GET("/posts/:id", func(c *gin.Context) {
		db, err := sql.Open("mysql", "root:polarbear21*@tcp(localhost:3306)/winter")
		defer db.Close()

		// Check for database connection errors
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to connect to database"})
			return
		}

		// Get the ID parameter from the request
		idParam := c.Param("id")

		// Convert the ID parameter to an integer
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid ID"})
			return
		}

		// Query the database for a post with the matching ID
		var post Post
		row := db.QueryRow("SELECT * FROM liketable WHERE ID=?", id)
		err = row.Scan(&post.ID, &post.Title, &post.Content, &post.Like)

		// Check for query errors
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "post not found"})
				return
			}
			c.JSON(500, gin.H{"error": "failed to retrieve post"})
			return
		}

		// Increment the like count for the post and update the database
		post.Like++
		_, err = db.Exec("UPDATE liketable SET `Like` = ? WHERE ID =?", post.Like, post.ID)

		// Printing error message.
		if err != nil {
			log.Println("Error:", err)
		}

		// Check for update errors
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to update post"})
			return
		}

		// Return the post to the client
		c.JSON(200, post)
	})

	router.Run()
}
