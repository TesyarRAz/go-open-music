package model

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type Playlist struct {
	ID        int       `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"insertedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    int
	User      *User   `json:"creator"`
	Users     []*User `json:"users" gorm:"many2many:playlist_users"`
	Songs     []*Song `json:"songs" gorm:"many2many:playlist_songs"`
}

func (p *Playlist) AfterCreate(db *gorm.DB) error {
	return db.Model(p).Association("Users").Append(&User{ID: p.UserID})
}

func (p *Playlist) BeforeDelete(db *gorm.DB) error {
	var err error
	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		if e := db.Model(p).Association("Users").Clear(); e != nil {
			err = e
		}
	}()

	go func() {
		defer wg.Done()
		if e := db.Model(p).Association("Songs").Clear(); e != nil {
			err = e
		}
	}()

	wg.Add(2)
	wg.Wait()

	return err
}
