package Forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func CreateComment(db *sql.DB, content string, userID int, topicID int) error {
	log.Printf("Attempting to add comment: content=%s, userID=%d, topicID=%d\n", content, userID, topicID)

	_, err := db.Exec("INSERT INTO Commentaires (content, user_id, topic_id) VALUES (?, ?, ?)", content, userID, topicID)
	if err != nil {
		log.Printf("erreur lors de l'insertion du commentaire : %s\n", err)
		return fmt.Errorf("erreur lors de l'insertion du commentaire: %w", err)
	}

	log.Printf("Commentaire ajouté avec succès !\n")
	return nil
}

func AddMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("STATIC/HTML/comment.html") 
		if err != nil {
			log.Printf("Template execution: %s", err)
			http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
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
	userIDStr := r.Form.Get("user_id") 
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Identifiant utilisateur invalide", http.StatusBadRequest)
		return
	}

	topicID, err := strconv.Atoi(r.Form.Get("topic_id"))
	if err != nil {
		http.Error(w, "Identifiant sujet invalide", http.StatusBadRequest)
		return
	}

	statusCode, db := Open()
	if statusCode != 0 {
		http.Error(w, "Erreur lors de la connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close() 

	err = CreateComment(db, content, userID, topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Commentaire ajouté avec succès"))
}

