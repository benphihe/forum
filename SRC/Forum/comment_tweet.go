package Forum

import (
  "html/template"
  "log"
  "net/http"
  "strconv"
  "time"
)


func CommentTweet(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("STATIC/HTML/commenttweet.html")
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

  idTweet, err := strconv.Atoi(r.Form.Get("id_tweet"))
  if err != nil {
    http.Error(w, "Invalid tweet ID", http.StatusBadRequest)
    return
  }

  content := r.Form.Get("content")
  userID := r.Form.Get("id_user")

  if len(content) > 250 {
    http.Error(w, "Le contenu du commentaire ne peut pas dépasser 250 caractères", http.StatusBadRequest)
    return
  }

  currentTime := time.Now().Format("2006-01-02 15:04:05")

  _, err = db.Exec("INSERT INTO comment_tweet (id_tweet, content, id_user, created_at) VALUES (?, ?, ?, ?)", idTweet, content, userID, currentTime)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusCreated)
  w.Write([]byte("Commentaire créé avec succès"))
}