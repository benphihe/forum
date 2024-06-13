package Forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
)


func CreateComment(postID int, userID int, content string) (int, error) {
	log.Printf("Lancement de CreateComment avec : postID=%d, userID=%d, content=%s\n", postID, userID, content)

	result, err := db.Exec("INSERT INTO commentaire_post (id_post, id_user, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {
		log.Printf("Erreur lors de l'insertion du commentaire : %s\n", err)
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
	log.Println("AddComment handler called")

	err := r.ParseForm()
	if err != nil {
		log.Printf("Erreur lors de l'analyse du formulaire: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postIDStr := r.Form.Get("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Printf("Invalid post ID: %s\n", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Form.Get("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user ID: %s\n", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	content := r.Form.Get("content")
	if content == "" {
		log.Println("Content is empty")
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	log.Printf("Received comment: postID=%d, userID=%d, content=%s\n", postID, userID, content)

	_, err = CreateComment(postID, userID, content)
	if err != nil {
		log.Printf("Erreur lors de la création du commentaire: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}


