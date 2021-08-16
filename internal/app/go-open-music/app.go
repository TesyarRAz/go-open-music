package goopenmusic

import (
	"github.com/TesyarRAz/go-open-music/internal/app/go-open-music/collab"
	"github.com/TesyarRAz/go-open-music/internal/app/go-open-music/playlist"
	"github.com/TesyarRAz/go-open-music/internal/app/go-open-music/song"
	"github.com/TesyarRAz/go-open-music/internal/app/go-open-music/upload"
	"github.com/TesyarRAz/go-open-music/internal/app/go-open-music/user"
	"github.com/TesyarRAz/go-open-music/internal/pkg/config"
	"github.com/TesyarRAz/go-open-music/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func setupQueue() (*amqp.Connection, *amqp.Channel) {
	conn := config.NewQueue()
	ch, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}

	return conn, ch
}

func Run() {
	r := gin.Default()

	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	config.Init()

	db := config.NewDatabase()
	queue, ch := setupQueue()
	cache := config.NewCache()

	defer queue.Close()
	defer ch.Close()

	songController := song.SongController{Db: db}
	userController := user.UserController{Db: db}
	playlistController := playlist.NewController(db, ch, cache)
	collabController := collab.CollabController{Db: db}
	uploadController := upload.UploadController{}

	authMiddleware := middleware.AuthMiddleware{Db: db}

	r.Static("/static", "static")

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
		auth.POST("/exports/playlists/:playlistId", playlistController.Export)

		auth.POST("/collaborations", collabController.Store)
		auth.DELETE("/collaborations", collabController.Destroy)
	}

	r.POST("upload/pictures", uploadController.Store)

	r.Run(":5000")
}
