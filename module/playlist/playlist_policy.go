package playlist

import "github.com/TesyarRAz/go-open-music/model"

type PlaylistPolicy struct {
	Playlist *model.Playlist
}

func (p *PlaylistPolicy) CanAccess(user *model.User) bool {
	for _, collaborator := range p.Playlist.Users {
		if collaborator.ID == user.ID {
			return true
		}
	}

	return false
}
