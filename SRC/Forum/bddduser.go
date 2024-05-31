package Forum

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func Open() (int, *sql.DB) {
	db, err := sql.Open("sqlite3", "BDD/Users.db")
	if err != nil {
		fmt.Println(err)
		fmt.Print("error ouvertur base")
		return 500, db
	}
	return 0, db
}

func CreateUser(pseudo string, password string, email string) error {
	log.Printf("CreateUser a reçu : pseudo=%s, password=%s, email=%s\n", pseudo, password, email)
	_, err := db.Exec("insert into Utilisateurs (pseudo, password, email) values (?, ?, ?)", pseudo, password, email)
	if err != nil {
		log.Printf("Erreur lors de l'insertion des données : %s\n", err)
	} else {
		log.Printf("Send envoie : pseudo=%s, password=%s, email=%s\n", pseudo, password, email)
	}
	return err
}
