package model

import "time"

type Playlist struct {
	ID        int       `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	UserName  string    `json:"username"`
	CreatedAt time.Time `json:"insertedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	User      User      `json:"user" gorm:"foreignKey:UserName"`
	Songs     []Song    `json:"songs" gorm:"many2many:playlist_songs"`
}
