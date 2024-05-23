package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbname := "main.db"
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Connected to the SQLite database:", dbname)

	// Créer une nouvelle base de données
	db2, err := sql.Open("sqlite3", "ma_base_de_donnees.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db2.Close()

	signup := func() {
		_, err := db2.Exec("CREATE TABLE IF NOT EXISTS utilisateurs (id INTEGER PRIMARY KEY, email TEXT UNIQUE, password TEXT, fullname TEXT)")
		if err != nil {
			log.Fatal(err)
		}
	}
	signup()

	rows, err := db2.Query("SELECT * FROM utilisateurs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id   int
		email string
		password string
		fullname string
	)
	for rows.Next() {
		err := rows.Scan(&id, &email, &password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, email, password)
	}
}


func register() {
	http.HandleFunc("/register", inscriptionHandler)
	http.ListenAndServe(":5000", nil)
}

func inscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	fullname := r.FormValue("fullname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	db, err := sql.Open("sqlite3", "./ma_base_de_donnees")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO utilisateurs (id, email, password, fullname) VALUES (NULL, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, password, fullname)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Utilisateur ajouté à la base de données.")

	http.ServeFile(w, r, "./front/index.html")

	rows, err := db.Query("SELECT * FROM utilisateurs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id int
		email string
		password string
		fullname string
	)
	for rows.Next() {
		err := rows.Scan(&id, &em, &mp, &nm)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, em, mp, nm)
	}
}