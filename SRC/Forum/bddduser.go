package Forum

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func Open() (int, *sql.DB) {
	db, err := sql.Open("sqlite3", "BDD/Users.db") //lancer depuis : (../../bdd.go) lancer depuis serveur.go : (./BDD/ProjetForum.db) le chemin du projet devra changer dependant de l'endroit exectution
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
	}
	return err
}

func Send(pseudo string, password string, email string) {
	log.Printf("Send envoie : pseudo=%s, password=%s, email=%s\n", pseudo, password, email)
	err := CreateUser(pseudo, password, email)
	if err != nil {
		log.Fatal(err)
	}
}
