package playlist

import (
	"net/http"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlaylistController struct {
	Db *gorm.DB
}

func (p *PlaylistController) Index(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	var playlists []model.Playlist

	if err := p.Db.Model(&user).Association("playlists").Find(&playlists); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    indexResponse(playlists),
		"message": "Berhasil membuat playlist",
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

	if err := p.Db.Model(&user).Association("playlists").Append(&playlist); err != nil {
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

func (p *PlaylistController) StoreSong(c *gin.Context) {
	playlistId := c.Param("playlistId")

	var (
		request  StoreSongRequest
		playlist model.Playlist
		song     model.Song
	)

	if err := p.Db.First(&playlist, "id = ?", playlistId).Error; err != nil {
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

	if err := p.Db.First(&song, "id = ?", request.SongId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := p.Db.Model(&playlist).Association("songs").Append(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berhasil membuat playlist",
	})
}
