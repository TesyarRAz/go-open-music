package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int            `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;notNull" binding:"required"`
	Password     string         `json:"password" gorm:"notNull" binding:"required"`
	FullName     string         `json:"fullname" gorm:"notNull" binding:"required"`
	AccessToken  sql.NullString `json:"accessToken" gorm:"default:null"`
	RefreshToken sql.NullString `json:"refreshToken" gorm:"default:null"`
	CreatedAt    time.Time      `json:"insertedAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	Playlists    []*Playlist    `json:"playlists"`
}
