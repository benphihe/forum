package main

import (
    "fmt"
    "net/http"
    "test/SRC/Forum"
    "text/template"
    "log"

)

func main() {

    http.HandleFunc("/", HomeHandler)
    http.HandleFunc("/connexion", Forum.Connec)
    http.HandleFunc("/inscription", Forum.InscriptionPage)
    fs := http.FileServer(http.Dir("STATIC"))
    http.Handle("/STATIC/", http.StripPrefix("/STATIC/", fs))

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