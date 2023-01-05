package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Set up the login route and template
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})

	// Handle form submission
	r.POST("/login", func(c *gin.Context) {
		// Get form values
		username := c.PostForm("username")
		password := c.PostForm("password")

		c.String(http.StatusOK, "Welcome "+username+"!\n")
		c.String(http.StatusOK, "your password is ,", password)
		// Check login credentials
		if isValidLogin(username, password) {
			// Set session cookie and redirect
			c.SetCookie("session_id", "12345", 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusFound, "/dashboard")
		} else {
			// Show error message
			c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"error": "Invalid login credentials"}) // <- HTML 파일이 없어도 login.tmpl의 이름으로 만든다.
		}
	})

	// Protected route
	r.GET("/dashboard", func(c *gin.Context) {
		// Check for valid session cookie
		if _, err := c.Cookie("session_id"); err == nil {
			c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{})
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
	})

	r.Run()
}

func isValidLogin(username, password string) bool {
	// Replace with actual login check
	return username == "user" && password == "pass"
}
