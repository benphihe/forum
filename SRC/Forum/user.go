package Forum

import (
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserInfoFromDB(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := map[string]interface{}{
		"Email":   user.Email,
		"id_user": user.ID,
		"pseudo":  user.Name,
	}

	tmpl, err := template.ParseFiles("STATIC/HTML/user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
