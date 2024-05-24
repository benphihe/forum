package main

import (
	"net/http"
	
	"Forum"
	
)

func main() {

	http.HandleFunc("/", Forum.AccueilPage)
	http.HandleFunc("/connexion", Forum.ConnexionPage)
	http.HandleFunc("/inscription", Forum.InscriptionPage)
	http.HandleFunc("/post", Forum.PostPage)
	http.HandleFunc("/comment", Forum.CommentPage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// http.HandleFunc("/sent", Accueil.GetSign)
	// http.HandleFunc("/sentConnect", Accueil.GetSignConnect)
	// http.HandleFunc("/sentText", Accueil.GetPostInformation)
	// http.HandleFunc("/sentCmt", Accueil.GetCmtInformation)

	http.ListenAndServe(":8080", nil)

}
