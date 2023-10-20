package models

type UserFullInfo struct {
	ID            int64
	Username      string
	FirstName     string
	LastName      string
	SubManageType string
	Role          string
	Favorites     []Anime
}

type Anime struct {
	Name string
	URL  string
}
