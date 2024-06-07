package Forum

import (
	"fmt"
	"html/template"
	"log"
	"strconv"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func CreatePost(pseudo string, id_user int, content_post string) (int, error) {
	log.Printf("Lancement de CreatePost avec : pseudo=%s, id_user=%d, content_post=%s\n", pseudo, id_user, content_post)

	_, db := Open()
	if db == nil {
		return 0, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	log.Printf("CreatePost a reçu : pseudo=%s, id_user=%d, content_post=%s\n", pseudo, id_user, content_post)
	result, err := db.Exec("INSERT INTO Post (pseudo, id_user, content_post) VALUES (?, ?, ?)", pseudo, id_user, content_post)
	if err != nil {
		log.Printf("erreur lors de l'insertion des données : %s\n", err)
		return 0, err
	} else {
		id, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		log.Printf("Post créé avec succès : pseudo=%s, id_user=%d, content_post=%s, id_post=%d\n", pseudo, id_user, content_post, id)
		return int(id), nil
	}
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

	content_post := r.Form.Get("content_post")
	id_user := globalUserID
	pseudo := globalPseudo

	log.Printf("Received form data: content_post=%s, id_user=%d, pseudo=%s\n", content_post, id_user, pseudo)

	statusCode, db := Open()
	if statusCode != 0 {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = CreatePost(pseudo, id_user, content_post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	action := r.Form.Get("action")
	if action == "like" {
		postIDStr := r.Form.Get("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			log.Printf("Invalid post ID: %s\n", postIDStr)
			return
		}

		err = AddLike(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
		return
	}
}

