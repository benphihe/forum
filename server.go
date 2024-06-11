package main

import (
	"fmt"
	"net/http"
	"test/SRC/Forum"
)

func main() {
	fs := http.FileServer(http.Dir("STATIC"))
	http.Handle("/STATIC/", http.StripPrefix("/STATIC/", fs))

	Forum.Open()
	http.HandleFunc("/", AuthMiddleware(Forum.DisplayPostsFrombdd))
	http.HandleFunc("/user", AuthMiddleware(Forum.DisplayUserInfo))
	http.HandleFunc("/user", AuthMiddleware(Forum.UserHandler))
	http.HandleFunc("/post/", AuthMiddleware(Forum.DisplayPost))
	http.HandleFunc("/inscription", Forum.InscriptionPage)
	http.HandleFunc("/connexion", Forum.Connexion)
	http.HandleFunc("/addpost", AuthMiddleware(Forum.AddPost))
	http.HandleFunc("/search", AuthMiddleware(Forum.SearchPosts))

	http.ListenAndServe(":8080", nil)
	fmt.Println("Server Start in localhost:8080")
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Forum.DisplayPost(w, r)
	} else if r.Method == http.MethodPost {
		Forum.AddComment(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

type key int

const (
	UserKey key = iota
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/inscription", http.StatusSeeOther)
			return
		}

		isValid, err := Forum.IsUUIDInDB(cookie.Value)
		if err != nil || !isValid {
			http.Redirect(w, r, "/inscription", http.StatusSeeOther)
			return
		}

		user, err := Forum.GetUserInfoFromDB(r)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		next(w, r.WithContext(ctx))
	}
}
