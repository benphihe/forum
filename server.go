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
	http.HandleFunc("/user", Forum.UserHandler)
	http.HandleFunc("/comment", Forum.AddComment)
	http.HandleFunc("/inscription", Forum.InscriptionPage)
	http.HandleFunc("/post", Forum.AddPost)
	http.HandleFunc("/tweet", Forum.AddTweet)
	http.HandleFunc("/comment_tweet", Forum.CommentTweet)
	http.HandleFunc("/connexion", Forum.Connexion)

	http.ListenAndServe(":8080", nil)
	fmt.Println("Server Start in localhost:8080")
}
