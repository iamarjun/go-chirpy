package database

type Chirp struct {
	ID       int    `json:"id"`
	Chirp    string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type Token struct {
	IsRevoked bool   `json:"is_revoked"`
	Timestamp string `json:"timestamp"`
}

type ResponseUser struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type ResponseUserWithToken struct {
	ResponseUser
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func UserToResponseUser(user User) ResponseUser {
	return ResponseUser{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
}
func UserToResponseUserWithToken(user User) ResponseUserWithToken {
	return ResponseUserWithToken{
		ResponseUser: UserToResponseUser(user),
	}
}
