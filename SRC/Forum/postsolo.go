package Forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := GetPost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, post)
}

func GetPost(postID int) (map[string]string, error) {
	_, db := Open()
	if db == nil {
		return nil, fmt.Errorf("erreur d'ouverture de la base de données")
	}

	row := db.QueryRow("SELECT id_post, id_user, pseudo, content_post FROM post WHERE id_post = ?", postID)

	var id_post, id_user, pseudo, content_post string
	if err := row.Scan(&id_post, &id_user, &pseudo, &content_post); err != nil {
		return nil, err
	}

	post := map[string]string{
		"id_post":      id_post,
		"id_user":      id_user,
		"pseudo":       pseudo,
		"content_post": content_post,
	}

	return post, nil
}

func DisplayPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := GetPost(postID)
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
		t.Execute(w, struct {
			Post     map[string]string
			Comments []map[string]string
		}{
			Post: post,
		})
		return
	}
}

func CreateComment(postID int, userID int, content string) (int, error) {
	_, db := Open()
	if db == nil {
		return 0, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO commentaire_post (content, id_user, id_post) VALUES (?, ?, ?)", content, userID, postID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	user, err := GetUserInfoFromDB(r)
	if err != nil {
		http.Error(w, "Error getting user info", http.StatusInternalServerError)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	_, err = CreateComment(postID, user.ID, content)
	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}
