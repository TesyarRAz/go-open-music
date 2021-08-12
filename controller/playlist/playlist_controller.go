package playlist

import (
	"net/http"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/TesyarRAz/go-open-music/policy"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlaylistController struct {
	Db *gorm.DB
}

func (p *PlaylistController) Index(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	playlists, err := userPlaylists(p.Db, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"playlists": indexResponse(playlists),
		},
		"message": "Berhasil mengambil playlist",
	})
}

func (p *PlaylistController) Store(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	var playlist model.Playlist

	if err := c.ShouldBindJSON(&playlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	if err := p.Db.Model(&user).Association("Playlists").Append(&playlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    storeResponse(playlist),
		"message": "Berhasil membuat playlist",
	})
}

func (p *PlaylistController) Destroy(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	playlistId := c.Param("playlistId")

	var playlist model.Playlist

	if err := p.Db.Preload("Users").First(&playlist, "id = ?", playlistId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Mengecek hanya user pembuat yang bisa menghapus
	if user.ID != playlist.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "playlist orang ini woe",
		})
		return
	}

	if err := p.Db.Delete(&playlist).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil menghapus playlist",
	})
}

func (p *PlaylistController) StoreSong(c *gin.Context) {
	playlistId := c.Param("playlistId")
	user := c.MustGet("user").(*model.User)

	var (
		request  storeSongRequest
		playlist model.Playlist
		song     model.Song
	)

	if err := p.Db.Preload("Users").First(&playlist, "id = ?", playlistId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	policy := policy.PlaylistPolicy{Playlist: &playlist}

	// Mengecek apakah user itu bisa mengakses playlistnya
	if !policy.CanAccess(user) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "playlist orang ini woe",
		})
		return
	}

	if err := p.Db.First(&song, "id = ?", request.SongId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := p.Db.Model(&playlist).Association("Songs").Append(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berhasil memasukan lagu ke playlist",
		"data":    storeResponse(playlist),
	})
}

func (p *PlaylistController) ShowSong(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	playlistId := c.Param("playlistId")

	var playlist model.Playlist

	if err := p.Db.Preload("Songs").Preload("Users").First(&playlist, "id = ?", playlistId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	policy := policy.PlaylistPolicy{Playlist: &playlist}

	// Mengecek apakah user itu bisa mengakses playlistnya
	if !policy.CanAccess(user) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "playlist orang ini woe",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil playlist",
		"data": gin.H{
			"songs": showResponse(playlist),
		},
	})
}

func (p *PlaylistController) DestroySong(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	playlistId := c.Param("playlistId")

	var (
		request  destroySongRequest
		playlist model.Playlist
		song     model.Song
	)

	if err := p.Db.Preload("Users").First(&playlist, "id = ?", playlistId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	policy := policy.PlaylistPolicy{Playlist: &playlist}

	// Mengecek apakah user itu bisa mengakses playlistnya
	if !policy.CanAccess(user) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "playlist orang ini woe",
		})
		return
	}

	if err := p.Db.First(&song, "id = ?", request.SongId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := p.Db.Model(&playlist).Association("Songs").Delete(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil playlist",
	})
}
