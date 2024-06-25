package Forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// GetUUIDFromCookie tente de récupérer l'UUID stocké dans un cookie de session d'une requête HTTP.
func GetUUIDFromCookie(r *http.Request) (string, error) {
	// Tente de récupérer le cookie "session_token" de la requête.
	cookie, err := r.Cookie("session_token")
	// Gère les erreurs liées à la récupération du cookie.
	if err != nil {
		// Si aucun cookie n'est trouvé, renvoie une erreur spécifique.
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("aucun cookie 'session_token' trouvé")
		}
		// Pour toute autre erreur, la renvoie telle quelle.
		return "", err
	}

	// Extrait l'UUID du cookie.
	uuid := cookie.Value

	// Vérifie si l'UUID est présent dans la base de données.
	isInDB, err := IsUUIDInDB(uuid)
	// Gère les erreurs potentielles lors de la vérification dans la base de données.
	if err != nil {
		return "", fmt.Errorf("erreur lors de la vérification de l'UUID dans la base de données : %s", err)
	}

	// Si l'UUID n'est pas trouvé dans la base de données, renvoie une erreur.
	if !isInDB {
		return "", fmt.Errorf("l'UUID du cookie n'est pas dans la base de données")
	}

	// Si tout est en ordre, renvoie l'UUID.
	return uuid, nil
}

// IsUUIDInDB vérifie si un UUID donné est présent dans la base de données.
func IsUUIDInDB(uuid string) (bool, error) {
	// Tente d'ouvrir une connexion à la base de données.
	_, db := Open()
	// Si la connexion échoue, renvoie une erreur.
	if db == nil {
		return false, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	// Assure que la connexion à la base de données sera fermée à la fin de la fonction.
	defer db.Close()

	// Exécute une requête pour récupérer tous les UUID des utilisateurs.
	rows, err := db.Query("SELECT UUID FROM Utilisateurs")
	// Gère les erreurs potentielles de la requête.
	if err != nil {
		return false, err
	}
	// Assure que le curseur des résultats sera fermé à la fin de la fonction.
	defer rows.Close()

	// Parcourt les résultats de la requête.
	var dbUUID sql.NullString
	for rows.Next() {
		// Tente de lire chaque ligne de résultat.
		err = rows.Scan(&dbUUID)
		// Gère les erreurs potentielles lors de la lecture.
		if err != nil {
			log.Printf("Erreur lors de la lecture des résultats : %s", err)
			return false, err
		}
		// Vérifie si l'UUID correspond à celui recherché.
		if dbUUID.Valid {
			if dbUUID.String == uuid {
				// Si un correspondance est trouvée, renvoie vrai.
				return true, nil
			}
		}
	}

	// Si aucune correspondance n'est trouvée, renvoie faux.
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

