package Forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func GetUUIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("aucun cookie 'session_token' trouvé")
		}
		return "", err
	}

	uuid := cookie.Value

	isInDB, err := IsUUIDInDB(uuid)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la vérification de l'UUID dans la base de données : %s", err)
	}

	if !isInDB {
		return "", fmt.Errorf("l'UUID du cookie n'est pas dans la base de données")
	}

	return uuid, nil
}

func IsUUIDInDB(uuid string) (bool, error) {
	_, db := Open()
	if db == nil {
		return false, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	rows, err := db.Query("SELECT UUID FROM Utilisateurs")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var dbUUID sql.NullString
	for rows.Next() {
		err = rows.Scan(&dbUUID)
		if err != nil {
			log.Printf("Erreur lors de la lecture des résultats : %s", err)
			return false, err
		}
		if dbUUID.Valid {
			if dbUUID.String == uuid {
				return true, nil
			}
		}
	}

	return false, nil
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    uuid, err := GetUUIDFromCookie(r)
    if err != nil {
        http.Redirect(w, r, "/connexion", http.StatusSeeOther)
        return
    }
    fmt.Println("UUID récupéré du cookie avec succès : ", uuid)

    w.WriteHeader(http.StatusOK)
}


