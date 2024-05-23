package Forum

import (
    _ "database/sql"
    "net/http"
    "encoding/json"
)

type Message struct {
    ID       int    `json:"id"`
    Content  string `json:"content"`
    UserID   int    `json:"user_id"`
    TopicID  int    `json:"topic_id"`
}

func AddMessage(w http.ResponseWriter, r *http.Request) {
    var m Message
    err := json.NewDecoder(r.Body).Decode(&m)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.Exec("INSERT INTO messages (content, user_id, topic_id) VALUES (?, ?, ?)", m.Content, m.UserID, m.TopicID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Message ajouté avec succès"))
}
