package Forum

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

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

		Send(pseudo, password, email)

		fmt.Fprintf(w, "Inscription réussie\n")
		fmt.Fprintf(w, "Pseudo: %s\n", pseudo)
		fmt.Fprintf(w, "Email: %s\n", email)
	}
}
