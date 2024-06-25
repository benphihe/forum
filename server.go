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
	http.HandleFunc("/like", AuthMiddleware(Forum.HandleLike))
	http.HandleFunc("/logout", AuthMiddleware(Forum.Logout))
	http.HandleFunc("/comment/", AuthMiddleware(Forum.AddComment))

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

// AuthMiddleware est une fonction qui prend en paramètre un gestionnaire HTTP (http.HandlerFunc) et retourne un autre gestionnaire HTTP.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// La fonction retournée est une fonction anonyme qui prend un ResponseWriter et une Request comme paramètres.
	return func(w http.ResponseWriter, r *http.Request) {
		// Tentative de récupération du cookie "session_token" de la requête.
		cookie, err := r.Cookie("session_token")
		// Si le cookie n'est pas trouvé (ou une autre erreur survient), redirige l'utilisateur vers la page d'inscription.
		if err != nil {
			http.Redirect(w, r, "/inscription", http.StatusSeeOther)
			return // Arrête l'exécution de la fonction ici.
		}

		// Vérifie si la valeur du cookie (UUID) est présente dans la base de données.
		isValid, err := Forum.IsUUIDInDB(cookie.Value)
		// Si une erreur survient lors de la vérification ou si l'UUID n'est pas valide, redirige vers la page d'inscription.
		if err != nil || !isValid {
			http.Redirect(w, r, "/inscription", http.StatusSeeOther)
			return // Arrête l'exécution de la fonction ici.
		}

		// Si le cookie est valide et présent dans la base de données, exécute le gestionnaire HTTP suivant dans la chaîne.
		next(w, r)
	}
}