package Forum

import (
	"database/sql"
	"fmt"
	"net/http"
)

type User struct {
	ID    int    `db:"id_user"`
	Name  string `db:"pseudo"`
	Email string `db:"email"`
}

func GetUserInfoFromDB(r *http.Request) (*User, error) {
	uuid, err := GetUUIDFromCookie(r)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'UUID du cookie : %s", err)
	}

	_, db := Open()
	if db == nil {
		return nil, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	query := `SELECT id_user, pseudo, email FROM Utilisateurs WHERE UUID = ?`
	row := db.QueryRow(query, uuid)

	user := &User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("aucun utilisateur trouvé avec l'UUID : %s", uuid)
		}
		return nil, fmt.Errorf("erreur lors de la récupération des informations de l'utilisateur : %s", err)
	}

	return user, nil
}
