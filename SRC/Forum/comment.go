package Forum

import (
	"html/template"
	"log"
	"net/http"
)

type Message struct {
	ID      int
	Content string
	UserID  int
	TopicID int
}

func AddMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("STATIC/HTML/comment.html")
		if err != nil {
			log.Fatalf("Template execution: %s", err)
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
	userID := r.Form.Get("user_id")
	topicID := r.Form.Get("topic_id")

	_, err = db.Exec("INSERT INTO messages (content, user_id, topic_id) VALUES (?, ?, ?)", content, userID, topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message ajouté avec succès"))
}
