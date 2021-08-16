package playlist

import (
	"context"

	"github.com/TesyarRAz/go-open-music/internal/pkg/model"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func userPlaylists(db *gorm.DB, cache *redis.Client, c context.Context, user *model.User) ([]model.Playlist, error) {
	var playlists []model.Playlist
	var playlistUsers []model.PlaylistUser

	// db.Joins("playlists_users ON playlists_users.playlist_id = playlists.id").Where("playlist_users.user_id = ?", user.ID).Preload("User").Find(&playlists)
	err := db.Preload("Playlist.User").Find(&playlistUsers, "user_id = ?", user.ID).Error

	for _, p := range playlistUsers {
		// Load Cache
		if ok := p.Playlist.LoadCacheSongs(cache, c); !ok {
			if err := db.Model(&p.Playlist).Association("Songs").Find(&p.Playlist.Songs); err != nil {
				return nil, err
			}
		}
		playlists = append(playlists, *p.Playlist)
	}

	return playlists, err
}
