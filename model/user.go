package model

import "time"

type User struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	Username     string     `json:"username" gorm:"index:,unique,notNull" binding:"required"`
	Password     string     `json:"password" gorm:"notNull" binding:"required"`
	FullName     string     `json:"fullname" gorm:"notNull" binding:"required"`
	AccessToken  string     `json:"accessToken"`
	RefreshToken string     `json:"refreshToken"`
	CreatedAt    time.Time  `json:"insertedAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	Playlists    []Playlist `json:"playlists" gorm:"foreignKey:UserName"`
}
