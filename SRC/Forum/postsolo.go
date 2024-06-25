package Forum

import (
	"database/sql"
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

	row := db.QueryRow("SELECT id_user, pseudo, content_post FROM post WHERE id_post = ?", postID)

	var id_user, pseudo, content_post string
	if err := row.Scan(&id_user, &pseudo, &content_post); err != nil {
		return nil, err
	}

	post := map[string]string{
		"id_post":      fmt.Sprintf("%d", postID),
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
		t.Execute(w, struct {
			Post     map[string]string
			Comments []map[string]string
		}{
			Post:     post,
			Comments: comments,
		})
		return
	}
}
func CreateComment(userID int, content string, postID int) (map[string]interface{}, error) {
	_, db := Open()
	if db == nil {
		return nil, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	_, err := GetPost(postID)
	if err != nil {
		return nil, fmt.Errorf("le post avec l'ID %d n'existe pas", postID)
	}

	var pseudo string
	err = db.QueryRow("SELECT pseudo FROM Utilisateurs WHERE id_user = ?", userID).Scan(&pseudo)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du pseudo de l'utilisateur: %v", err)
	}

	res, err := db.Exec("INSERT INTO commentaire_post (content, id_user, id_post, pseudo) VALUES (?, ?, ?, ?)", content, userID, postID, pseudo)
	if err != nil {
		return nil, err
	}

	commentID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	comment := map[string]interface{}{
		"id_commentaire_post": commentID,
		"content":             content,
		"id_user":             userID,
		"id_post":             postID,
		"pseudo":              pseudo,
	}

	return comment, nil
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := strings.TrimPrefix(r.URL.Path, "/comment/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "ID du post invalide", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Le contenu du commentaire ne peut pas être vide", http.StatusBadRequest)
		return
	}

	user, err := GetUserInfoFromDB(r)
	if err != nil {
		http.Error(w, "Impossible d'identifier l'utilisateur", http.StatusInternalServerError)
		return
	}

	userID := user.ID

	comment, err := CreateComment(userID, content, postID)
	if err != nil {
		http.Error(w, "Erreur lors de la création du commentaire: "+err.Error(), http.StatusInternalServerError)
		return
	}

	success, err := AddCommentToDB(postID, userID, comment["content"].(string), db)
	if err != nil || !success {
		http.Error(w, "Erreur lors de l'ajout du commentaire à la base de données: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Commentaire ajouté avec succès"))
}

func AddCommentToDB(postID int, userID int, content string, db *sql.DB) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("la connexion à la base de données n'est pas initialisée")
	}
	if content == "" {
		return false, fmt.Errorf("le contenu du commentaire ne peut pas être vide")
	}

	_, err := db.Exec("INSERT INTO commentaire_post (id_post, id_user, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {
		return false, fmt.Errorf("erreur lors de l'ajout du commentaire : %s", err)
	}
	return true, nil
}

func GetComments(postID int) ([]map[string]string, error) {
	_, db := Open()
	if db == nil {
		return nil, fmt.Errorf("erreur d'ouverture de la base de données")
	}
	defer db.Close()

	rows, err := db.Query("SELECT pseudo, content FROM commentaire_post WHERE id_post = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []map[string]string
	for rows.Next() {
		var pseudo sql.NullString
		var content sql.NullString
		if err := rows.Scan(&pseudo, &content); err != nil {
			return nil, err
		}
		comment := map[string]string{
			"pseudo":          pseudo.String,
			"content_comment": content.String,
		}
		if !pseudo.Valid {
			comment["pseudo"] = "Anonyme"
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
