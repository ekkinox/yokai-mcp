package domain

type Book struct {
	ID       int32  `json:"id"`
	Title    string `json:"title"`
	Genre    string `json:"genre"`
	Synopsis string `json:"synopsis"`
}
