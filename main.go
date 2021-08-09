package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Song struct {
	ID        int       `json:"id" gorm:"primarykey"`
	Title     string    `json:"title" gorm:"notNull" binding:"required"`
	Year      int       `json:"year" gorm:"notNull" binding:"required"`
	Performer string    `json:"performer" gorm:"notNull" binding:"required"`
	Genre     string    `json:"genre" gorm:"notNull" binding:"required"`
	Duration  string    `json:"duration" gorm:"notNull" binding:"required"`
	CreatedAt time.Time `json:"insertedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func main() {
	r := gin.Default()

	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))))
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&Song{})

	r.POST("/songs", func(c *gin.Context) {
		var song Song

		if err := c.ShouldBindJSON(&song); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})

			return
		}

		if err := db.Create(&song).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})

			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Berhasil membuat lagu",
			"data": gin.H{
				"songId": strconv.Itoa(song.ID),
			},
		})
	})

	r.GET("/songs", func(c *gin.Context) {
		var songs []Song

		if err := db.Find(&songs).Error; err != nil {
			log.Panic(err.Error())
		}

		var resources []gin.H

		for _, s := range songs {
			resources = append(resources, gin.H{
				"id":        strconv.Itoa(s.ID),
				"title":     s.Title,
				"performer": s.Performer,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"songs": resources,
			},
		})
	})

	r.GET("/songs/:songId", func(c *gin.Context) {
		songId := c.Param("songId")

		var song Song

		if err := db.First(&song, "id = ?", songId).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"song": struct {
					ID        string    `json:"id"`
					Title     string    `json:"title" binding:"required"`
					Year      int       `json:"year" binding:"required"`
					Performer string    `json:"performer" binding:"required"`
					Genre     string    `json:"genre" binding:"required"`
					Duration  int       `json:"duration" binding:"required"`
					CreatedAt time.Time `json:"insertedAt"`
					UpdatedAt time.Time `json:"updatedAt"`
				}{
					ID:        strconv.Itoa(song.ID),
					Title:     song.Title,
					Year:      song.Year,
					Performer: song.Performer,
					Genre:     song.Genre,
					Duration: func() int {
						duration, _ := strconv.Atoi(song.Duration)

						return duration
					}(),
					CreatedAt: song.CreatedAt,
					UpdatedAt: song.UpdatedAt,
				},
			},
		})
	})

	r.PUT("/songs/:songId", func(c *gin.Context) {
		songId := c.Param("songId")

		var (
			song    Song
			newSong struct {
				ID        int       `json:"id"`
				Title     string    `json:"title" binding:"required"`
				Year      string    `json:"year" binding:"required"`
				Performer string    `json:"performer" binding:"required"`
				Genre     string    `json:"genre" binding:"required"`
				Duration  string    `json:"duration" binding:"required"`
				CreatedAt time.Time `json:"insertedAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			}
		)

		if err := db.First(&song, "id = ?", songId).Error; err != nil {
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

		if err := db.Model(&song).Updates(&newSong).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Berhasil merubah lagu",
			"data": struct {
				ID        int       `json:"id"`
				Title     string    `json:"title" binding:"required"`
				Year      string    `json:"year" binding:"required"`
				Performer string    `json:"performer" binding:"required"`
				Genre     string    `json:"genre" binding:"required"`
				Duration  int       `json:"duration" binding:"required"`
				CreatedAt time.Time `json:"insertedAt"`
				UpdatedAt time.Time `json:"updatedAt"`
			}{
				ID:        newSong.ID,
				Title:     newSong.Title,
				Year:      newSong.Year,
				Performer: newSong.Performer,
				Genre:     newSong.Genre,
				Duration: func() int {
					duration, _ := strconv.Atoi(newSong.Duration)

					return duration
				}(),
				CreatedAt: newSong.CreatedAt,
				UpdatedAt: newSong.UpdatedAt,
			},
		})
	})

	r.DELETE("/songs/:songId", func(c *gin.Context) {
		songId := c.Param("songId")

		var song Song
		if err := db.First(&song, "id = ?", songId).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}

		if err := db.Delete(&song).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Berhasil menghapus lagu",
			"data":    song,
		})
	})

	r.Run(":5000")
}
