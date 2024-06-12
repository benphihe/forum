package Forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetUserIDFromSession(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}

	sessionToken := cookie.Value

	_, db := Open()
	if db == nil {
		return ""
	}
	defer db.Close()

	var userID string
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID)
	if err != nil {
		return ""
	}

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

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idPost := r.FormValue("id_post")
		userId := GetUserIDFromSession(r)

		if idPost == "" || userId == "" {
			http.Error(w, "Invalid post ID or user ID", http.StatusBadRequest)
			return
		}

		_, db := Open()
		defer db.Close()

		var ownerID string
		err := db.QueryRow("SELECT id_user FROM post WHERE id_post = ?", idPost).Scan(&ownerID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if ownerID != userId {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = db.Exec("DELETE FROM post WHERE id_post = ?", idPost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}





