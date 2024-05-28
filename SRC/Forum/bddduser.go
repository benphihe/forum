package Forum

import (
	"database/sql"
	"log"
)

type User struct {
	ID       int    `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"mdp"`
}

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "BDD/users.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	create table if not exists users (id integer not null primary key, pseudo text, password text, email text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
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



