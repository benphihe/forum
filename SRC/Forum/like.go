package Forum

import (
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func AddLike(id_post int, id_user int) error {
	_, db := Open()
	if db == nil {
		return fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	_, err := db.Exec("INSERT INTO likes (id_post, id_user) VALUES (?, ?)", id_post, id_user)
	if err != nil {
		return err
	}
	return nil
}

func GetLikeCount(id_post int) (int, error) {
	_, db := Open()
	if db == nil {
		return 0, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	var likeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE id_post = ?", id_post).Scan(&likeCount)
	if err != nil {
		return 0, err
	}
	return likeCount, nil
}

func HandleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id_post, err := strconv.Atoi(r.Form.Get("postID"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	id_user_str := GetUserIDFromSession(r)
	if id_user_str == "" {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	id_user, err := strconv.Atoi(id_user_str)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = AddLike(id_post, id_user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	likeCount, err := GetLikeCount(id_post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%d", likeCount)
}


