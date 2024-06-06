package Forum

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func CreateComment(postID int, userID int, content string,) (int, error) {
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
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("STATIC/HTML/comment.html")
		if err != nil {
			log.Printf("Template execution: %s", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	content := r.Form.Get("content")
	postIDStr := r.Form.Get("post_id")
	userIDStr := r.Form.Get("user_id")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		log.Printf("Invalid post ID: %s\n", postIDStr)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Printf("Invalid user ID: %s\n", userIDStr)
		return
	}

	log.Printf("Received form data: postID=%d, content=%s, userID=%d\n", postID, content, userID)

	_, err = CreateComment(postID, userID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}





