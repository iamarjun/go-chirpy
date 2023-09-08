package database

type Chirp struct {
	ID    int    `json:"id"`
	Chirp string `json:"body"`
}
