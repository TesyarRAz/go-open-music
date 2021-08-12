package model

import "time"

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
