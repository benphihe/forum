package Forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"database/sql"
)

func SearchPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var id_category int64
	switch query {
	case "#Mylife":
		id_category = 1
	case "#Cinema":
		id_category = 2
	case "#Sport":
		id_category = 3
	default:
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	posts, err := GetPostsByCategory(id_category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("STATIC/HTML/search_results.html")
	if err != nil {
		log.Printf("Template execution: %s", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, posts)
}

func GetPostsByCategory(categoryID int64) ([]map[string]interface{}, error) {
	_, db := Open()
	if db == nil {
		return nil, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	rows, err := db.Query("SELECT id_post, id_user, pseudo, content_post, id_category FROM post WHERE id_category = ?", categoryID)
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

		post := map[string]interface{}{
			"id_post":      id_post,
			"id_user":      id_user,
			"pseudo":       pseudo,
			"content_post": content_post,
			"category_name": categoryName,
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
