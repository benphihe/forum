package Forum

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var globalUserID int
var globalPseudo string

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

		log.Println("Email: ", email)
		log.Println("Password: ", password)

		authenticated, uuid, err := Authenticate(email, password)
		if err != nil {
			if errors.Is(err, errors.New("invalid email or password")) {
				http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
				return
			}
			log.Fatal(err)
		}

		if userID != 0 {
			globalUserID = userID
			globalPseudo = pseudo
			log.Println("Connexion réussie")
			log.Println("UUID: ", uuid)
			http.Redirect(w, r, "http://localhost:8080/user?email="+url.QueryEscape(email), http.StatusSeeOther)
		} else {
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		}
	}
}


func Authenticate(email string, password string) (bool, string, error) {
	_, db := Open()
	if db == nil {
		return false, "", fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	var dbPassword string
	var uuid string
	err := db.QueryRow("SELECT password, uuid FROM Utilisateurs WHERE email = ?", email).Scan(&dbPassword, &uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}

	log.Println("Email from DB: ", email)
	log.Println("Password from DB: ", dbPassword)
	log.Println("UUID from DB: ", uuid)

	err = VerifyHash(dbPassword, password)
	if err != nil {
		return false, "", err
	}

	_, _, _, err = GetUser(email)
	if err != nil {
		return false, "", err
	}

	return true, uuid, nil
}

func VerifyHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}



