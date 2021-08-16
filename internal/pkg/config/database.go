package config

import (
	"fmt"

	"github.com/TesyarRAz/go-open-music/internal/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	var (
		config = &gorm.Config{
			FullSaveAssociations: false,
		}
		db  *gorm.DB
		err error
	)

	if AppConfig.DB_DRIVER == "pgsql" {
		db, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", AppConfig.DB_HOST, AppConfig.DB_USERNAME, AppConfig.DB_PASSWORD, AppConfig.DB_NAME, AppConfig.DB_PORT)), config)
	} else if AppConfig.DB_DRIVER == "sqlite" {
		db, err = gorm.Open(sqlite.Open(AppConfig.DB_HOST))
	}
	if err != nil {
		panic(err.Error())
	}

	if d, err := db.DB(); err == nil {
		d.SetMaxIdleConns(10)
		d.SetMaxOpenConns(100)
	}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Song{})
	db.AutoMigrate(&model.Playlist{})
	db.AutoMigrate(&model.PlaylistUser{})

	if err := db.SetupJoinTable(&model.Playlist{}, "Users", &model.PlaylistUser{}); err != nil {
		panic(err.Error())
	}

	if err := db.SetupJoinTable(&model.User{}, "AllPlaylists", &model.PlaylistUser{}); err != nil {
		panic(err.Error())
	}

	return db
}
