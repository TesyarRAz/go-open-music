package song

import "time"

type UpdateRequest struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Year      string    `json:"year" binding:"required"`
	Performer string    `json:"performer" binding:"required"`
	Genre     string    `json:"genre" binding:"required"`
	Duration  string    `json:"duration" binding:"required"`
	CreatedAt time.Time `json:"insertedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
