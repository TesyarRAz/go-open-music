package song

import (
	"log"
	"net/http"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SongController struct {
	Db *gorm.DB
}

func (s *SongController) Index(c *gin.Context) {
	var songs []model.Song

	if err := s.Db.Find(&songs).Error; err != nil {
		log.Panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"songs": indexResponse(songs),
		},
	})
}

func (s *SongController) Store(c *gin.Context) {
	var song model.Song

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	if err := s.Db.Create(&song).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berhasil membuat lagu",
		"data":    storeResponse(song),
	})
}

func (s *SongController) Show(c *gin.Context) {
	songId := c.Param("songId")

	var song model.Song

	if err := s.Db.First(&song, "id = ?", songId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"song": showResponse(song),
		},
	})
}

func (s *SongController) Update(c *gin.Context) {
	songId := c.Param("songId")

	var (
		song    model.Song
		newSong updateRequest
	)

	if err := s.Db.First(&song, "id = ?", songId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&newSong); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := s.Db.Model(&song).Updates(newSong).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil merubah lagu",
		"data":    updateResponse(newSong),
	})
}

func (s *SongController) Destroy(c *gin.Context) {
	songId := c.Param("songId")

	var song model.Song
	if err := s.Db.First(&song, "id = ?", songId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if err := s.Db.Delete(&song).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil menghapus lagu",
	})
}
