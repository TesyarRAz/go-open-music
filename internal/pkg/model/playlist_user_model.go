package model

type PlaylistUser struct {
	UserID     int
	PlaylistID int
	User       *User
	Playlist   *Playlist
}
