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

type ResponseUserWithToken struct {
	ResponseUser
	Token string `json:"token"`
}

func UserToResponseUser(user User) ResponseUser {
	return ResponseUser{
		ID:    user.ID,
		Email: user.Email,
	}
}
func UserToResponseUserWithToken(user User) ResponseUserWithToken {
	return ResponseUserWithToken{
		ResponseUser: UserToResponseUser(user),
	}
}
