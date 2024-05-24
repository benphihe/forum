package Forum

import (
	"log"
	"net/http"
)

func Websitestart() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "style.css") })
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	log.Printf("Server running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "STATIC/HTML/Accueil.html")
}
