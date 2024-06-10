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
    http.HandleFunc("/post/", PostHandler)
    http.HandleFunc("/inscription", Forum.InscriptionPage)
    http.HandleFunc("/connexion", Forum.Connexion)
    http.HandleFunc("/addpost", Forum.AddPost)
    // http.HandleFunc("/cookies", Forum.SignOutHandler)

    http.ListenAndServe(":8080", nil)
    fmt.Println("Server Start in localhost:8080")
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        Forum.DisplayPost(w, r)
    } else if r.Method == http.MethodPost {
        Forum.AddComment(w, r)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

