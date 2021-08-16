package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
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

type ExportPlaylist struct {
	TargetEmail string `json:"targetEmail" binding:"required,email"`
	Playlist    *Playlist
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

func (p *Playlist) SaveCacheSongs(rdb *redis.Client, c context.Context) {
	rdb.LRem(c, fmt.Sprintf("playlist:%d", p.ID), 0, nil)

	var wg sync.WaitGroup
	for _, s := range p.Songs {
		wg.Add(1)
		go func(s Song) {
			defer wg.Done()
			str, err := json.Marshal(s)
			if err != nil {
				log.Fatalln(err.Error())
				return
			}

			rdb.LPush(c, fmt.Sprintf("playlist:%d:songs", p.ID), str)
		}(*s)
	}

	wg.Wait()
}

func (p *Playlist) LoadCacheSongs(rdb *redis.Client, c context.Context) bool {
	slice := rdb.LRange(c, fmt.Sprintf("playlist:%d:songs", p.ID), 0, -1)

	if len(slice.Val()) == 0 {
		return false
	}

	for _, s := range slice.Val() {
		var song Song
		err := json.Unmarshal([]byte(s), &song)
		if err != nil {
			log.Fatalln(err.Error())
			continue
		}

		p.Songs = append(p.Songs, &song)
	}

	return true
}
