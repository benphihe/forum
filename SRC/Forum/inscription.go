package Forum

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

func InscriptionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var u User
    err := json.NewDecoder(r.Body).Decode(&u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    stmt, err := db.Prepare("INSERT INTO utilisateurs(fullname, email, mdp) VALUES(?,?,?)")
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

    response := map[string]interface{}{
        "id":      id,
        "message": "Inscription réussie!",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


