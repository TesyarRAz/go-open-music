package main

import (
	"github.com/TesyarRAz/go-open-music/config"
	"github.com/TesyarRAz/go-open-music/controller/playlist"
	"github.com/TesyarRAz/go-open-music/controller/song"
	"github.com/TesyarRAz/go-open-music/controller/user"
	"github.com/TesyarRAz/go-open-music/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	config.Init()

	db := config.NewDatabase()

	songController := song.SongController{Db: db}
	userController := user.UserController{Db: db}
	playlistController := playlist.PlaylistController{Db: db}

	authMiddleware := middleware.AuthMiddleware{Db: db}

	r.GET("/songs", songController.Index)
	r.GET("/songs/:songId", songController.Show)
	r.POST("/songs", songController.Store)
	r.PUT("/songs/:songId", songController.Update)
	r.DELETE("/songs/:songId", songController.Destroy)

	r.POST("/users", userController.Store)
	r.POST("/authentications", userController.Login)
	r.PUT("/authentications", userController.Refresh)

	auth := r.Group("/", authMiddleware.Auth)
	{
		auth.GET("/playlists", playlistController.Index)
		auth.POST("/playlists", playlistController.Store)
		auth.POST("/playlists/{playlistId}/songs", playlistController.StoreSong)
	}

	r.Run(":5000")
}
