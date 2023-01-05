package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

// Post represents a single post on the post board.
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
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
		result, err := db.Exec("INSERT INTO postsave (Title, Content) VALUES (?, ?)", post.Title, post.Content)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		// Get the ID of the new post
		id, err := result.LastInsertId()
		//	fmt.Println(id)
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
		idParam := c.Param("id")
		fmt.Println("Received ID:", idParam)
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid ID"})
			return
		}
		for _, post := range Posts {
			if post.ID == id {
				c.JSON(200, post)
				return
			}
		}
		c.JSON(404, gin.H{"error": "post not found"})
	})

	router.Run()
}
