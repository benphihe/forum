package main

import (
	"fmt"
	"net/http"
	"test/SRC/Forum"
)

func main() {
	fs := http.FileServer(http.Dir("STATIC"))
	http.Handle("/STATIC/", http.StripPrefix("/STATIC/", fs))

	Forum.Open()
	http.HandleFunc("/", Forum.DisplayPosts)
	http.HandleFunc("/inscription", Forum.InscriptionPage)
	http.HandleFunc("/connexion", Forum.Connexion)
	http.HandleFunc("/cookies", Forum.SignOutHandler)

	http.Handle("/user", Forum.AuthMiddleware(http.HandlerFunc(Forum.UserHandler)))
	http.Handle("/comment", Forum.AuthMiddleware(http.HandlerFunc(Forum.AddComment)))
	http.Handle("/post", Forum.AuthMiddleware(http.HandlerFunc(Forum.AddPost)))
	http.Handle("/tweet", Forum.AuthMiddleware(http.HandlerFunc(Forum.AddTweet)))
	http.Handle("/comment_tweet", Forum.AuthMiddleware(http.HandlerFunc(Forum.CommentTweet)))

	http.ListenAndServe(":8080", nil)
	fmt.Println("Server Start in localhost:8080")
}

