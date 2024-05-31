package Forum

import (
	"fmt"
	"log"
	"text/template"
	"net/http"
	"errors"
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

		authenticated, err := Authenticate(email, password)
        if err != nil {
            if errors.Is(err, errors.New("Invalid email or password")) {
                http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
                return
            }
            log.Fatal(err)
        }

        if authenticated {
            fmt.Fprintf(w, "Connexion réussie")
        } else {
            http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
        }
    }
}


func Authenticate(email string, password string) (bool, error) {
    var storedPassword string
    err := db.QueryRow("SELECT password FROM Utilisateurs WHERE email = ?", email).Scan(&storedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }
    return password == storedPassword, nil
}


