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

func CreatePost(db *sql.DB, content string, userID int, pseudo string) error {
	log.Printf("Attempting to add post: content=%s, userID=%d, pseudo=%s\n", content, userID, pseudo)

	_, err := db.Exec("INSERT INTO Post (content_post, id_user, pseudo) VALUES (?, ?, ?)", content, userID, pseudo)
	if err != nil {
		log.Printf("Error inserting post: %s\n", err)
		return fmt.Errorf("Error inserting post: %w", err)
	}

	log.Printf("Post added successfully!\n")
	return nil
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("STATIC/HTML/post.html")
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

	content := r.Form.Get("content_post")
	userIDStr := r.URL.Query().Get("user_id")
	pseudo := r.URL.Query().Get("pseudo")

	log.Printf("Received form data: content=%s, userIDStr=%s, pseudo=%s\n", content, userIDStr, pseudo)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Printf("Invalid user ID: %s\n", userIDStr)
		return
	}

	statusCode, db := Open()
	if statusCode != 0 {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = CreatePost(db, content, userID, pseudo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Post added: content=%s, userID=%d, pseudo=%s\n", content, userID, pseudo)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post added successfully"))
}






