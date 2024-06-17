package Forum

import (
    "net/http"
    "time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
    cookie := http.Cookie{
        Name:    "session_token",
        Value:   "",
        Expires: time.Now().Add(-1 * time.Hour),
        Path:    "/",
    }
    http.SetCookie(w, &cookie)

    http.Redirect(w, r, "/connexion", http.StatusSeeOther)
}
