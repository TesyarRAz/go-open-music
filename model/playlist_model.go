package model

import (
	"time"

	"gorm.io/gorm"
)

type Playlist struct {
	ID        int       `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"insertedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    int
	User      *User   `json:"user"`
	Songs     []*Song `json:"songs" gorm:"many2many:playlist_songs"`
}

func (p *Playlist) BeforeDelete(db *gorm.DB) error {
	return db.Model(p).Association("Songs").Clear()
}
