package Forum

import (
	"database/sql"
	"fmt"
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
