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
	fmt.Println("Base de données ouverte avec succès")
	defer db.Close()

	rows, err := db.Query("SELECT UUID FROM Utilisateurs")
	if err != nil {
		log.Printf("Erreur lors de l'exécution de la requête SQL : %s", err)
		return false, err
	}
	fmt.Println("Requête SQL exécutée avec succès")
	defer rows.Close()

	fmt.Println("UUID du cookie : ", uuid)

	var dbUUID sql.NullString
	var uuids []string
	for rows.Next() {
		err = rows.Scan(&dbUUID)
		if err != nil {
			log.Printf("Erreur lors de la lecture des résultats : %s", err)
			return false, err
		}
		if dbUUID.Valid {
			uuids = append(uuids, dbUUID.String)
			if dbUUID.String == uuid {
				fmt.Println("UUID correspondant trouvé dans la base de données : ", dbUUID.String)
				fmt.Println("Le cookie est égal à l'UUID dans la base de données")
				return true, nil
			} else {
				fmt.Println("Le cookie n'est pas égal à l'UUID dans la base de données")
			}
		}
	}

	fmt.Println("UUIDs dans la base de données : ", uuids)

	return false, nil
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	uuid, err := GetUUIDFromCookie(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println("UUID récupéré du cookie avec succès")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("UUID: " + uuid))
	fmt.Println("Réponse envoyée avec succès")
}

