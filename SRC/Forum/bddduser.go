package Forum

import (
	"database/sql"
	"log"
)

var db *sql.DB

func Open() {
	var err error
	db, err = sql.Open("sqlite3", "BDD/users.db")
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(pseudo string, password string, email string) error {
	_, err := db.Exec("insert into users (pseudo, password, email) values (?, ?, ?)", pseudo, password, email)
	return err
}

func Send(pseudo string, password string, email string) {
	err := CreateUser(pseudo, password, email)
	if err != nil {
		log.Fatal(err)
	}
}



