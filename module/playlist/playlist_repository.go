package playlist

import (
	"github.com/TesyarRAz/go-open-music/model"
	"gorm.io/gorm"
)

func userPlaylists(db *gorm.DB, user *model.User) ([]model.Playlist, error) {
	var playlists []model.Playlist
	var playlistUsers []model.PlaylistUser

	// db.Joins("playlists_users ON playlists_users.playlist_id = playlists.id").Where("playlist_users.user_id = ?", user.ID).Preload("User").Find(&playlists)
	err := db.Preload("Playlist.User").Find(&playlistUsers, "user_id = ?", user.ID).Error

	for _, p := range playlistUsers {
		playlists = append(playlists, *p.Playlist)
	}

	return playlists, err
}
