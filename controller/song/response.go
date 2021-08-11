package song

import (
	"strconv"
	"time"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/gin-gonic/gin"
)

// Response Converter
func indexResponse(songs []model.Song) interface{} {
	var resources []gin.H

	for _, s := range songs {
		resources = append(resources, gin.H{
			"id":        strconv.Itoa(s.ID),
			"title":     s.Title,
			"performer": s.Performer,
		})
	}

	return resources
}

func storeResponse(song model.Song) interface{} {
	return gin.H{
		"songId": strconv.Itoa(song.ID),
	}
}

func showResponse(song model.Song) interface{} {
	return struct {
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
	}
}

func updateResponse(song UpdateRequest) interface{} {
	return struct {
		ID        int       `json:"id"`
		Title     string    `json:"title" binding:"required"`
		Year      string    `json:"year" binding:"required"`
		Performer string    `json:"performer" binding:"required"`
		Genre     string    `json:"genre" binding:"required"`
		Duration  int       `json:"duration" binding:"required"`
		CreatedAt time.Time `json:"insertedAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}{
		ID:        song.ID,
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
	}
}
