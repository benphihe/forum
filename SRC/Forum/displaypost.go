package Forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func GetUserIDFromSession(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("No session cookie found:", err)
		return ""
	}

	uuid := cookie.Value
	fmt.Println("UUID from cookie:", uuid)

	_, db := Open()
	if db == nil {
		fmt.Println("Failed to open database")
		return ""
	}
	defer db.Close()

	var userID string
	err = db.QueryRow("SELECT id_user FROM Utilisateurs WHERE UUID = ?", uuid).Scan(&userID)
	if err != nil {
		fmt.Println("Failed to retrieve user ID from UUID:", err)
		return ""
	}
	fmt.Println("Retrieved user ID:", userID)

	return userID
}


func DisplayPostsFrombdd(w http.ResponseWriter, r *http.Request) {
	id_category := r.URL.Query().Get("category")
	userID := GetUserIDFromSession(r)
	posts, err := GetPosts(userID, id_category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("STATIC/HTML/acceuil.html")
		if err != nil {
			log.Printf("Template execution: %s", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		t.Execute(w, posts)
		return
	}
}

func GetPosts(userID string, id_category string) ([]map[string]interface{}, error) {
	_, db := Open()
	if db == nil {
		return nil, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	var rows *sql.Rows
	var err error
	if id_category != "" {
		rows, err = db.Query("SELECT id_post, id_user, pseudo, content_post, id_category FROM post WHERE id_category = ?", id_category)
	} else {
		rows, err = db.Query("SELECT id_post, id_user, pseudo, content_post, id_category FROM post")
	}

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




