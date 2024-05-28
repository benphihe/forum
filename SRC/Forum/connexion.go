package Forum

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

type User struct {
    ID       int
    FullName string
    Email    string
    Password string
}

func ConnexionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var u User
    err := json.NewDecoder(r.Body).Decode(&u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    row := db.QueryRow("SELECT id, fullname, email, mdp FROM utilisateurs WHERE email = ? AND mdp = ?", u.Email, u.Password)
    err = row.Scan(&u.ID, &u.FullName, &u.Email, &u.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Connexion réussie"))
}


