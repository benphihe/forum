package main

import (
	"fmt"
	"net/http"
	"test/SRC/Forum"

)

func main() {

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/connexion", Forum.Connec)
	http.HandleFunc("/inscription", Forum.InscriptionPage)
	fs := http.FileServer(http.Dir("STATIC"))
	http.Handle("/STATIC/", http.StripPrefix("/STATIC/", fs))

	http.ListenAndServe(":8080", nil)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "STATIC/HTML/Acceuil.html")
	fmt.Println("Server is running")

}
