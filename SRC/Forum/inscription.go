package Forum

import (
	"fmt"
	//"log"
	"net/http"
	_"github.com/mattn/go-sqlite3"
    
    //envoie "../BDD"
)



func InscriptionPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	pseudo := r.FormValue("pseudo")
	password := r.FormValue("password")

	fmt.Println("Pseudo: ", pseudo)
	fmt.Println("Mot de passe: ", password)

	fmt.Fprintf(w, "Inscription réussie")
}

// func Send(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "send.html")
// 	err := r.ParseForm()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	pseudo := r.FormValue("Pseudo")
// 	password := r.FormValue("Password")

// 	fmt.Println(" Identidiant d'Inscription : ", pseudo, "/", password)
// 	http.Redirect(w, r, "/connexion", http.StatusSeeOther)
// 	statusenvoie, db := envoie.GestionData()
// 	status := envoie.NewUser(pseudo, password, db)
// 	if status == 0 && statusenvoie == 0 {
// 		fmt.Println("creation d'un nouveau utilisateur")
// 	} else {
// 		fmt.Println("echec de la creation d'un nouveau utilisateur")
// 	}
// }
