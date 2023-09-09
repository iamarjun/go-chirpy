package database

type Chirp struct {
	ID    int    `json:"id"`
	Chirp string `json:"body"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func UserToResponseUser(user User) ResponseUser {
	return ResponseUser{
		ID:    user.ID,
		Email: user.Email,
	}
}
