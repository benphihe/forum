package Forum

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func Open() {
	if _, err := os.Stat("BDD/users.db"); os.IsNotExist(err) {
		log.Fatal("Database file does not exist")
		return
	}

	var err error
	db, err = sql.Open("sqlite3", "./BDD/users.db")
	if err != nil {
		log.Fatal(err)
	}
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






