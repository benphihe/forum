package Forum

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

func InscriptionPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("STATIC/HTML/inscription.html")
		if err != nil {
			log.Fatalf("Template execution: %s", err)
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		pseudo := r.FormValue("pseudo")
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("passwordConfirm")
		email := r.FormValue("email")

		if password != passwordConfirm {
			http.Error(w, "Les mots de passe ne correspondent pas", http.StatusBadRequest)
			return
		}

		hashedPassword, err := Hash(password)
		if err != nil {
			http.Error(w, "Erreur lors du hachage du mot de passe", http.StatusInternalServerError)
			return
		}

		err = CreateUser(pseudo, hashedPassword, email, GenerateUUID())
		if err != nil {
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
	}
}

func GenerateUUID() string {
	return uuid.New().String()
}

func CreateUser(pseudo string, password string, email string, uuid string) error {
	_, db := Open()
	if db == nil {
		return fmt.Errorf("erreur d'ouverture de la base de données")
	}
	log.Printf("CreateUser a reçu : pseudo=%s, password=%s, email=%s, uuid=%s\n", pseudo, password, email, uuid)
	_, err := db.Exec("insert into Utilisateurs (pseudo, password, email, uuid) values (?, ?, ?, ?)", pseudo, password, email, uuid)
	if err != nil {
		log.Printf("erreur lors de l'insertion des données : %s\n", err)
	} else {
		log.Printf("Send envoie : pseudo=%s, password=%s, email=%s, uuid=%s\n", pseudo, password, email, uuid)
	}
	return err
}

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}


