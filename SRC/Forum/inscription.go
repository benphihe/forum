package Forum

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "log"
)


func Inscription() {
    var err error
    db, err = sql.Open("sqlite3", "ma_base_de_donnees.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    sqlStmt := `
    CREATE TABLE IF NOT EXISTS utilisateurs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        fullname TEXT,
        email TEXT UNIQUE,
        mdp TEXT
    );`
    _, err = db.Exec(sqlStmt)
    if err != nil {
        log.Printf("%q: %s\n", err, sqlStmt)
        return
    }

    http.HandleFunc("/inscription", inscriptionHandler)
}

func inscriptionHandler(w http.ResponseWriter, r *http.Request) {
    var u User
    err := json.NewDecoder(r.Body).Decode(&u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    stmt, err := db.Prepare("INSERT INTO utilisateurs(fullname, email, mdp) values(?,?,?)")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    res, err := stmt.Exec(u.FullName, u.Email, u.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    id, err := res.LastInsertId()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Inscription réussie, id=%d", id)
}