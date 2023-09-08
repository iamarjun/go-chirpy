package database

type Chirp struct {
	ID    int    `json:"id"`
	Chirp string `json:"body"`

}
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
