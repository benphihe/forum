package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"mdp"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "ma_base_de_donnees.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS utilisateurs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fullname TEXT,
		email TEXT UNIQUE,
		mdp TEXT
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	http.HandleFunc("/inscription", inscriptionHandler)
	http.HandleFunc("/connexion", connexionHandler)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {http.ServeFile(w, r, "style.css")})
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	log.Printf("Server running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func inscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO utilisateurs(fullname, email, mdp) values(?,?,?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(u.FullName, u.Email, u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Inscription réussie, id=%d", id)
}

func connexionHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM utilisateurs WHERE email = ? AND mdp = ?", u.Email, u.Password)
	err = row.Scan(&u.ID, &u.FullName, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprintf(w, "Connexion réussie, user=%+v", u)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "Accueil.html")
}
