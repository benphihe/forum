package Forum

type User struct {
	ID       int    `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"mdp"`
}