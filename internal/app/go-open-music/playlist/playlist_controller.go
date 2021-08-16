package playlist

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TesyarRAz/go-open-music/internal/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type PlaylistController struct {
	Db           *gorm.DB
	QueueChannel *amqp.Channel
	Queue        amqp.Queue
	Cache        *redis.Client
}

func NewController(db *gorm.DB, queueChannel *amqp.Channel, cache *redis.Client) *PlaylistController {
	queue, err := queueChannel.QueueDeclare(
		"playlist-queue", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)

	if err != nil {
		panic(err)
	}

	return &PlaylistController{
		Db:           db,
		QueueChannel: queueChannel,
		Queue:        queue,
		Cache:        cache,
	}
}

func (p *PlaylistController) Index(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	playlists, err := userPlaylists(p.Db, p.Cache, c, user)

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

	p.Cache.LRem(c, fmt.Sprintf("playlist:%d", playlist.ID), 0, -1)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil menghapus playlist",
	})
}

func (p *PlaylistController) Export(c *gin.Context) {
	var (
		playlistId = c.Param("playlistId")
		user       = c.MustGet("user").(*model.User)
		playlist   model.Playlist
		export     model.ExportPlaylist
	)

	if err := p.Db.Preload("Users").First(&playlist, "id = ?", playlistId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Load Cache
	if ok := playlist.LoadCacheSongs(p.Cache, c); !ok {
		if err := p.Db.Model(&playlist).Association("Songs").Find(&playlist.Songs); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}
	}

	if err := c.ShouldBindJSON(&export); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	export.Playlist = &playlist

	policy := PlaylistPolicy{Playlist: &playlist}

	// Mengecek apakah user itu bisa mengakses playlistnya
	if !policy.CanAccess(user) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "playlist orang ini woe",
		})
		return
	}

	data, err := json.Marshal(export)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := p.QueueChannel.Publish(
		"",           // exchange
		p.Queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berhasil mengexport playlist",
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

	policy := PlaylistPolicy{Playlist: &playlist}

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

	playlist.Songs = nil
	if err := p.Db.Model(&playlist).Association("Songs").Find(&playlist.Songs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}
	playlist.SaveCacheSongs(p.Cache, c)

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

	if err := p.Db.Preload("Users").First(&playlist, "id = ?", playlistId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Load Cache
	if ok := playlist.LoadCacheSongs(p.Cache, c); !ok {
		if err := p.Db.Model(&playlist).Association("Songs").Find(&playlist.Songs); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}
	}

	policy := PlaylistPolicy{Playlist: &playlist}

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

	policy := PlaylistPolicy{Playlist: &playlist}

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

	playlist.Songs = nil
	if err := p.Db.Model(&playlist).Association("Songs").Find(&playlist.Songs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}
	playlist.SaveCacheSongs(p.Cache, c)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mengambil playlist",
	})
}
