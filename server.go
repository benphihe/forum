package main

import (
    "net/http"
    "log"
    "test/SRC/Forum"
)

func main() {
    fs := http.FileServer(http.Dir("STATIC"))
    http.Handle("/STATIC/", http.StripPrefix("/STATIC/", fs))

    Forum.Websitestart()

    db, err := Forum.SetupDB()
    if err != nil {
        log.Fatal("Erreur de connexion à la base de données:", err)
    }

    http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
        Forum.ConnexionHandler(w, r, db)
    })
    http.HandleFunc("/inscription", func(w http.ResponseWriter, r *http.Request) {
        Forum.InscriptionHandler(w, r, db)
    })
    http.HandleFunc("/addMessage", func(w http.ResponseWriter, r *http.Request) {
        Forum.AddMessage(w, r, db)
    })

    log.Println("Serveur démarré sur le port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Erreur lors du démarrage du serveur:", err)
    }
}





