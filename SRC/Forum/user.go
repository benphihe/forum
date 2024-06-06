package Forum

import (
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func GetUser(email string) (string, int, string, error) {
	_, db := Open()
	if db == nil {
		return "", 0, "", fmt.Errorf("erreur d'ouverture de la base de données")
	}

	query := "SELECT email, id_user, pseudo FROM Utilisateurs WHERE email = ?"
	row := db.QueryRow(query, email)

	var userEmail string
	var userID int
	var username string

	err := row.Scan(&userEmail, &userID, &username)
	if err != nil {
		return "", 0, "", err
	}
	return userEmail, userID, username, nil
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	userEmail, userID, username, err := GetUser(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := map[string]interface{}{
		"Email":   userEmail,
		"id_user": userID,
		"pseudo":  username,
	}

	tmpl, err := template.ParseFiles("STATIC/HTML/user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

