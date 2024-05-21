package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

func Server() {
	http.HandleFunc("/", Home)
	fileServer := http.FileServer(http.Dir("templates/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets", fileServer))
	fmt.Println("http://localhost:8000/  Server is running in port 8000")
	http.ListenAndServe(":8000", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, error := template.ParseGlob("templates/*.html")
	if error != nil {
		panic(error)
	}
	tmpl.ExecuteTemplate(w, "index", nil)
}
