package goopenmusicqueue

import (
	"fmt"

	"github.com/TesyarRAz/go-open-music/internal/app/go-open-music-queue/playlist"
	"github.com/TesyarRAz/go-open-music/internal/pkg/config"
	"github.com/joho/godotenv"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	config.Init()

	conn := config.NewQueue()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	defer ch.Close()

	forever := make(chan bool)

	playlist.Handle(ch)

	fmt.Println("Close Using Ctrl + c")
	<-forever
}
