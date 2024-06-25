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
	http.HandleFunc("/user", AuthMiddleware(Forum.UserHandler))
	http.HandleFunc("/post/", AuthMiddleware(Forum.DisplayPost))
	http.HandleFunc("/inscription", Forum.InscriptionPage)
	http.HandleFunc("/connexion", Forum.Connexion)
	http.HandleFunc("/addpost", AuthMiddleware(Forum.AddPost))
	http.HandleFunc("/search", AuthMiddleware(Forum.SearchPosts))
	http.HandleFunc("/rules", Forum.DisplayRules)
	http.HandleFunc("/delete_post", Forum.DeletePost)
	http.HandleFunc("/addcomment", AuthMiddleware(Forum.AddComment))
	http.HandleFunc("/like", AuthMiddleware(Forum.HandleLike))
	http.HandleFunc("/logout", AuthMiddleware(Forum.Logout))

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

		next(w, r)
	}
}



