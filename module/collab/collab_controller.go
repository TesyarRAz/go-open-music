package collab

import (
	"net/http"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CollabController struct {
	Db *gorm.DB
}

func (co *CollabController) Store(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	var (
		request  storeRequest
		playlist model.Playlist
		collab   model.User
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	if err := co.Db.First(&playlist, "id = ?", request.PlaylistID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := co.Db.First(&collab, "id = ?", request.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Mengecek apakah user itu bisa mengakses playlistnya
	if user.ID != playlist.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "playlist orang ini woe",
		})
		return
	}

	if err := co.Db.Model(&playlist).Association("Users").Append(&collab); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berhasil masukan ke collab",
	})
}

func (co *CollabController) Destroy(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	var (
		request  destroyRequest
		playlist model.Playlist
		collab   model.User
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	if err := co.Db.First(&playlist, "id = ?", request.PlaylistID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := co.Db.First(&collab, "id = ?", request.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Mengecek apakah user itu bisa mengakses playlistnya
	if user.ID != playlist.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "playlist orang ini woe",
		})
		return
	}

	if err := co.Db.Model(&playlist).Association("Users").Delete(&collab); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berhasil hapus collab",
	})
}
