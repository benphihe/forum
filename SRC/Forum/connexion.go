package Forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

)

func connexion() {
	var err error
	db, err = sql.Open("sqlite3", "ma_base_de_donnees.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/connexion", connexionHandler)
}

func connexionHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM utilisateurs WHERE email = ? AND mdp = ?", u.Email, u.Password)
	err = row.Scan(&u.ID, &u.FullName, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprintf(w, "Connexion réussie, user=%+v", u)
}
