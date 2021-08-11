package playlist

import (
	"strconv"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/gin-gonic/gin"
)

func indexResponse(playlists []model.Playlist) interface{} {
	var resources []gin.H

	for _, s := range playlists {
		resources = append(resources, gin.H{
			"id":       strconv.Itoa(s.ID),
			"name":     s.Name,
			"username": s.UserName,
		})
	}

	return resources
}

func storeResponse(playlist model.Playlist) interface{} {
	return gin.H{
		"playlistId": strconv.Itoa(playlist.ID),
	}
}
