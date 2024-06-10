package Forum

import (
	"fmt"
	"html/template"
	"database/sql"
	"log"
	"strconv"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func CreatePost(pseudo string, id_user int, content_post string, id_category int) (int, error) {
	log.Printf("Lancement de CreatePost avec : pseudo=%s, id_user=%d, content_post=%s, id_category=%d\n", pseudo, id_user, content_post, id_category)

	_, db := Open()
	if db == nil {
		return 0, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	log.Printf("CreatePost a reçu : pseudo=%s, id_user=%d, content_post=%s, id_category=%d\n", pseudo, id_user, content_post, id_category)
	result, err := db.Exec("INSERT INTO Post (pseudo, id_user, content_post, id_category) VALUES (?, ?, ?, ?)", pseudo, id_user, content_post, id_category)
	if err != nil {
		log.Printf("erreur lors de l'insertion des données : %s\n", err)
		return 0, err
	} else {
		id, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		log.Printf("Post créé avec succès : pseudo=%s, id_user=%d, content_post=%s, id_category=%d, id_post=%d\n", pseudo, id_user, content_post, id_category, id)
		return int(id), nil 
	}
}

func AddPost(w http.ResponseWriter, r *http.Request) {
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

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	content_post := r.Form.Get("content_post")
	id_category, err := strconv.Atoi(r.Form.Get("category"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	id_user := globalUserID
	pseudo := globalPseudo

	log.Printf("Received form data: content_post=%s, id_user=%d, pseudo=%s, id_category=%d\n", content_post, id_user, pseudo, id_category)

	statusCode, db := Open()
	if statusCode != 0 {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = CreatePost(pseudo, id_user, content_post, id_category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DisplayPosts(w http.ResponseWriter, r *http.Request) {
	id_category := r.URL.Query().Get("category")
	var rows *sql.Rows
	var err error

	_, db := Open()
	if db == nil {
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if id_category == "" {
		rows, err = db.Query("SELECT pseudo, content_post, id_category FROM Post")
	} else {
		rows, err = db.Query("SELECT pseudo, content_post, id_category FROM Post WHERE id_category = ?", id_category)
	}

	if err != nil {
		log.Printf("Erreur lors de la récupération des données : %s\n", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Post struct {
		Pseudo      string
		ContentPost string
		CategoryID  sql.NullInt64
	}

	var posts []Post

	for rows.Next() {
		var pseudo, content_post string
		var categoryID sql.NullInt64
		err = rows.Scan(&pseudo, &content_post, &categoryID)
		if err != nil {
			log.Printf("Erreur lors de la lecture des données : %s\n", err)
			http.Error(w, "Error reading posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, Post{Pseudo: pseudo, ContentPost: content_post, CategoryID: categoryID})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Erreur lors de la finalisation de la lecture des données : %s\n", err)
		http.Error(w, "Error finalizing posts", http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("STATIC/HTML/acceuil.html")
	if err != nil {
		log.Printf("Template execution: %s", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, posts)
}


