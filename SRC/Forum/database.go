package Forum

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

func SetupDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "ma_base_de_donnees.db")
    if err != nil {
        log.Fatalf("Erreur lors de l'ouverture de la base de données: %v", err)
        return nil, err
    }

    return db, nil
}

