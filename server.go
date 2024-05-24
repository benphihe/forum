package main

import (
	"net/http"
	"test/SRC/Forum"
)

func main() {
	fs := http.FileServer(http.Dir("STATIC"))
	http.Handle("/STATIC/", http.StripPrefix("/STATIC/", fs))
	Forum.Websitestart()
	Forum.Connexion()
	Forum.Inscription()
	//http.HandleFunc("/post", Forum.PostPage)
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// http.HandleFunc("/sent", Accueil.GetSign)
	// http.HandleFunc("/sentConnect", Accueil.GetSignConnect)
	// http.HandleFunc("/sentText", Accueil.GetPostInformation)
	// http.HandleFunc("/sentCmt", Accueil.GetCmtInformation)

}
