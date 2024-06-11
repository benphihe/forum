package Forum

import (
	// "fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func DisplayUserInfo(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserInfoFromDB(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("STATIC/HTML/user.Html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
