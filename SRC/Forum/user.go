package Forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserInfoFromDB(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userPosts, err := GetUserPosts(strconv.Itoa(user.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := map[string]interface{}{
		"Email":   user.Email,
		"id_user": user.ID,
		"pseudo":  user.Name,
		"Content": userPosts,
	}

	tmpl, err := template.ParseFiles("STATIC/HTML/user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUserPosts(userID string) ([]map[string]interface{}, error) {
	_, db := Open()
	if db == nil {
		return nil, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	rows, err := db.Query("SELECT id_post, id_user, pseudo, content_post, id_category FROM post WHERE id_user = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryNames := map[int64]string{
		1: "#Mylife",
		2: "#Cinema",
		3: "#Sport",
	}

	var posts []map[string]interface{}
	for rows.Next() {
		var id_post, id_user, pseudo, content_post string
		var id_category sql.NullInt64
		if err := rows.Scan(&id_post, &id_user, &pseudo, &content_post, &id_category); err != nil {
			return nil, err
		}

		categoryName := "Unknown"
		if id_category.Valid {
			categoryName = categoryNames[id_category.Int64]
		}

		isOwner := (id_user == userID)

		post := map[string]interface{}{
			"id_post":       id_post,
			"id_user":       id_user,
			"pseudo":        pseudo,
			"content_post":  content_post,
			"category_name": categoryName,
			"is_owner":      isOwner,
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	_, db := Open()
	if db == nil {
		http.Error(w, fmt.Errorf("erreur d'ouverture de la base de données").Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err := db.Exec("DELETE FROM post WHERE id_post = ?", postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user", http.StatusSeeOther)
}
