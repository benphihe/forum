package Forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
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

	// Trier les posts par nombre de likes en ordre décroissant
	sort.Slice(posts, func(i, j int) bool {
		return posts[i]["like_count"].(int) > posts[j]["like_count"].(int)
	})

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

	query := `
	SELECT p.id_post, p.id_user, p.pseudo, p.content_post, p.id_category, 
	       (SELECT COUNT(*) FROM like WHERE like.id_post = p.id_post) as like_count
	FROM post p
	WHERE p.id_category = ?
	`
	rows, err := db.Query(query, categoryID)
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
		var like_count int
		if err := rows.Scan(&id_post, &id_user, &pseudo, &content_post, &id_category, &like_count); err != nil {
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
			"like_count":   like_count,
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

