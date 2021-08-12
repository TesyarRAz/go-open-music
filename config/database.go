package config

import (
	"fmt"

	"github.com/TesyarRAz/go-open-music/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", AppConfig.DB_HOST, AppConfig.DB_USERNAME, AppConfig.DB_PASSWORD, AppConfig.DB_NAME, AppConfig.DB_PORT)), &gorm.Config{
		FullSaveAssociations: false,
	})
	if err != nil {
		panic(err.Error())
	}

	db.Migrator().DropTable("playlist_songs")
	db.Migrator().DropTable(&model.Playlist{})
	db.Migrator().DropTable(&model.Song{})
	db.Migrator().DropTable(&model.User{})

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Song{})
	db.AutoMigrate(&model.Playlist{})

	return db
}
