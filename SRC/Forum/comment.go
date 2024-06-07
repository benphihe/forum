package Forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func CreateComment(postID int, userID int, content string) (int, error) {
	fmt.Println("Post ID:", postID)

	log.Printf("Lancement de CreateComment avec : postID=%d, userID=%d, content=%s\n", postID, userID, content)

	result, err := db.Exec("INSERT INTO commentaire_post (id_post, id_user, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {
		log.Printf("erreur lors de l'insertion du commentaire : %s\n", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Printf("Commentaire créé avec succès : postID=%d, userID=%d, content=%s, commentID=%d\n", postID, userID, content, id)
	return int(id), nil
}

func AddComment(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    parts := strings.Split(path, "/")
    if len(parts) < 3 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }

    postIDStr := parts[2]
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    fmt.Println("Post ID:", postID)

    if r.Method == "GET" {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        t, err := template.ParseFiles("STATIC/HTML/postsolo.html")
        if err != nil {
            log.Printf("Template execution: %s", err)
            http.Error(w, "Error executing template", http.StatusInternalServerError)
            return
        }
        t.Execute(w, nil)
        return
    }

    err = r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    userIDStr := r.Form.Get("userID")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    content := r.Form.Get("content")

    _, err = CreateComment(postID, userID, content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

func GetComments(postID int) ([]map[string]string, error) {
	_, db := Open()
	if db == nil {
		return nil, fmt.Errorf("erreur d'ouverture de la base de données")
	}

	rows, err := db.Query("SELECT id_commentaire_post, id_user, pseudo, content FROM commentaire_post WHERE id_post = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []map[string]string
	for rows.Next() {
		var id_comment, id_user, pseudo, content_comment string
		if err := rows.Scan(&id_comment, &id_user, &pseudo, &content_comment); err != nil {
			return nil, err
		}
		comment := map[string]string{
			"id_commentaire_post": id_comment,
			"id_user":             id_user,
			"pseudo":              pseudo,
			"content":             content_comment,
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func DisplayComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := GetComments(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("STATIC/HTML/postsolo.html")
		if err != nil {
			log.Printf("Template execution: %s", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		t.Execute(w, comments)
		return
	}
}

