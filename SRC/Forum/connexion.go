package Forum

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

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

			uuid, err := GetUUIDFromEmail(w, email)
			if err != nil {
				log.Printf("Erreur lors de la récupération de l'UUID : %s", err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
			log.Printf("UUID pour l'email %s est %s", email, uuid)

			http.Redirect(w, r, "/post", http.StatusSeeOther)
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

func GetUUIDFromEmail(w http.ResponseWriter, email string) (string, error) {
	_, db := Open()
	if db == nil {
		return "", fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	var uuid string
	err := db.QueryRow("SELECT uuid FROM Utilisateurs WHERE email = ?", email).Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération de l'UUID : %s", err)
	}

	cookie, err := CreateCookieWithUUID(uuid)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du cookie : %s", err)
	}

	http.SetCookie(w, cookie)

	return uuid, nil
}

func CreateCookieWithUUID(uuid string) (*http.Cookie, error) {
	if uuid == "" {
		return nil, fmt.Errorf("UUID vide")
	}

	expiration := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   uuid,
		Expires: expiration,
		Path:    "/",
	}
	return cookie, nil
}
