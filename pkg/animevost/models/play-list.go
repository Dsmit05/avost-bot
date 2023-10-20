package models

//go:generate easyjson -all
type Playlists []Playlist

//go:generate easyjson -all
type Playlist struct {
	SD      string `json:"std"`
	Preview string `json:"preview"`
	Name    string `json:"name"`
	HD      string `json:"hd"`
}
