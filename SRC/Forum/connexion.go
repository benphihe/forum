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

		userID, pseudo, err := AuthenticateAndGetUserID(email, password)
		if err != nil {
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
			return
		}

		if userID != 0 {
			globalUserID = userID
			globalPseudo = pseudo
			log.Println("Connexion réussie")

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		}
	}
}

func AuthenticateAndGetUserID(email string, password string) (int, string, error) {
	_, db := Open()
	if db == nil {
		return 0, "", fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	var dbPassword string
	var userID int
	var pseudo string
	err := db.QueryRow("SELECT id_user, password, pseudo FROM Utilisateurs WHERE email = ?", email).Scan(&userID, &dbPassword, &pseudo)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", errors.New("invalid email or password")
		}
		return 0, "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		return 0, "", errors.New("invalid email or password")
	}

	return userID, pseudo, nil
}

func VerifyHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
