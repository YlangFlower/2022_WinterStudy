package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

type loginInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login.html")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginData loginInfo
	err := json.NewDecoder(r.Body).Decode(&loginData)
	log.Println(err)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if loginData.Name == "" || loginData.Password == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Connect to database
	db, err := sql.Open("mysql", "root:polarbear21*@tcp(localhost:3306)/winter")
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Save login data to database
	_, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", loginData.Name, loginData.Password)
	if err != nil {
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Successfully saved to database")
}

func main() {
	// Open a new HTML window
	cmd := exec.Command("cmd", "/c", "start", "http://localhost:8080/login")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error opening HTML window:", err)
		return
	}

	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/login-submit", loginHandler)
	http.ListenAndServe(":8080", nil)
}
