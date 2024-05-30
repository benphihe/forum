package Forum

import (
	"fmt"
	"log"
	"text/template"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Connexion(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("STATIC/HTML/connexion.html")
		if err != nil {
			log.Fatalf("Template execution: %s", err)
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if Authenticate(email, password) {
			fmt.Fprintf(w, "Connexion réussie")
		} else {
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		}
	}
}


func Authenticate(email string, password string) bool {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return password == storedPassword
}


