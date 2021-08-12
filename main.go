package main

import (
	"github.com/TesyarRAz/go-open-music/config"
	"github.com/TesyarRAz/go-open-music/middleware"
	"github.com/TesyarRAz/go-open-music/module/collab"
	"github.com/TesyarRAz/go-open-music/module/playlist"
	"github.com/TesyarRAz/go-open-music/module/song"
	"github.com/TesyarRAz/go-open-music/module/user"
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
	collabController := collab.CollabController{Db: db}

	authMiddleware := middleware.AuthMiddleware{Db: db}

	r.GET("/songs", songController.Index)
	r.GET("/songs/:songId", songController.Show)
	r.POST("/songs", songController.Store)
	r.PUT("/songs/:songId", songController.Update)
	r.DELETE("/songs/:songId", songController.Destroy)

	r.POST("/users", userController.Store)
	r.POST("/authentications", userController.Login)
	r.PUT("/authentications", userController.Refresh)
	r.DELETE("/authentications", userController.DestroyToken)

	auth := r.Group("/", authMiddleware.Auth)
	{
		auth.GET("/playlists", playlistController.Index)
		auth.POST("/playlists", playlistController.Store)
		auth.DELETE("/playlists/:playlistId", playlistController.Destroy)
		auth.GET("/playlists/:playlistId/songs", playlistController.ShowSong)
		auth.POST("/playlists/:playlistId/songs", playlistController.StoreSong)
		auth.DELETE("/playlists/:playlistId/songs", playlistController.DestroySong)

		auth.POST("/collaborations", collabController.Store)
		auth.DELETE("/collaborations", collabController.Destroy)
	}

	r.Run(":5000")
}
