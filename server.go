package main

import (
	"fmt"
	"log"
	"net/http"
	"test/SRC/Forum"
	"text/template"
)

func main() {
	fs := http.FileServer(http.Dir("STATIC"))
	http.Handle("/STATIC/", http.StripPrefix("/STATIC/", fs))
	Forum.Open()
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/user", Forum.User)
	http.HandleFunc("/connexion", Forum.Connexion)
	http.HandleFunc("/comment", Forum.AddMessage)
	http.HandleFunc("/inscription", Forum.InscriptionPage)
	http.HandleFunc("/post", Forum.AddPost)
	http.HandleFunc("/tweet", Forum.AddTweet)
	http.HandleFunc("/comment_tweet", Forum.CommentTweet)

	http.ListenAndServe(":8080", nil)
	fmt.Println("Server Start in localhost:8080")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("STATIC/HTML/Accueil.html")
	if err != nil {
		log.Fatalf("Template execution: %s", err)
		return
	}
	t.Execute(w, nil)

}
