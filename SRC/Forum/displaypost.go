package Forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func GetPosts() ([]map[string]interface{}, error) {
    _, db := Open()
    if db == nil {
        return nil, fmt.Errorf("erreur d'ouverture de la base de données")
    }
    defer db.Close()

    rows, err := db.Query("SELECT id_post, id_user, id_category, pseudo, content_post FROM post")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []map[string]interface{}
    for rows.Next() {
        var id_post, id_user, pseudo, content_post string
        var id_category sql.NullInt64
        if err := rows.Scan(&id_post, &id_user, &id_category, &pseudo, &content_post); err != nil {
            return nil, err
        }

        post := map[string]interface{}{
            "id_post":      id_post,
            "id_user":      id_user,
            "id_category":  id_category.Int64,
            "pseudo":       pseudo,
            "content_post": content_post,
        }
        posts = append(posts, post)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return posts, nil
}


func DisplayPostsFromDatabase(w http.ResponseWriter, r *http.Request) {
	posts, err := GetPosts()
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



