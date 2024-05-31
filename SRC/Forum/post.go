package Forum

import (
  "html/template"
  "log"
  "net/http"
  "time"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t, err := template.ParseFiles("STATIC/HTML/index.html")
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

  pseudo := r.Form.Get("pseudo")
  content := r.Form.Get("content_post")
  userID := r.Form.Get("id_user")

  currentTime := time.Now().Format("2006-01-02 15:04:05")

  _, err = db.Exec("INSERT INTO posts (pseudo, content_post, id_user, created_at) VALUES (?, ?, ?, ?)", pseudo, content, userID, currentTime)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusCreated)
  w.Write([]byte("Post créé avec succès"))
}