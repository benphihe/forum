package Forum

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"

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

		log.Println("Email: ", email)
		log.Println("Password: ", password)

		authenticated, uuid, err := Authenticate(email, password, w, r)
		if err != nil {
			if errors.Is(err, errors.New("invalid email or password")) {
				http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
				return
			}
			log.Fatal(err)
		}

		if authenticated {
			if uuid != "" {
				log.Println("Connexion réussie")
				log.Println("UUID: ", uuid)
				http.Redirect(w, r, "http://localhost:8080/user?email="+url.QueryEscape(email), http.StatusSeeOther)
			} else {
				log.Println("UUID is empty")
			}
		} else {
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
		}
	}
}

func Authenticate(email string, password string, w http.ResponseWriter, r *http.Request) (bool, string, error) {
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

	// Si l'authentification est réussie, définissez le cookie et essayez de vous connecter automatiquement
	SetCookie(w, uuid)
	AutoLogin(w, r)

	return true, uuid, nil
}

func VerifyHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func SetCookie(w http.ResponseWriter, uuid string) {
	expirationTime := time.Now().Add(1 * time.Hour)
	cookie := &http.Cookie{
		Name:     "uuid",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		Expires:  expirationTime,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func AutoLogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("uuid")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("No cookie, user not logged in")
		} else {
			log.Println("Error getting cookie: ", err)
		}
		return
	}

	var email, password string
	err = db.QueryRow("SELECT email, password FROM Utilisateurs WHERE uuid = ?", cookie.Value).Scan(&email, &password)
	if err != nil {
		log.Println("Invalid UUID: ", err)
		return
	}

	isAuthenticated, _, err := Authenticate(email, password, w, r)
	if err != nil || !isAuthenticated {
		log.Println("Authentication failed: ", err)
		return
	}

	log.Println("User is logged in: ", email)

	http.Redirect(w, r, "/user/"+email, http.StatusSeeOther)
}
