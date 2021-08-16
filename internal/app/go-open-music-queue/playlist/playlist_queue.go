package playlist

import (
	"encoding/json"
	"log"
	"os"

	"github.com/TesyarRAz/go-open-music/internal/pkg/config"
	"github.com/TesyarRAz/go-open-music/internal/pkg/model"
	"github.com/streadway/amqp"
	"gopkg.in/gomail.v2"
)

func Handle(ch *amqp.Channel) {
	// Mengecek Apakah Ada Queue nya?

	q, err := ch.QueueDeclare(
		"playlist-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err.Error())
	}

	// Mengambil data jika ada
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		for d := range msgs {
			var export model.ExportPlaylist

			if err := json.Unmarshal(d.Body, &export); err != nil {
				log.Fatal(err.Error())
			}

			exportFile, err := os.CreateTemp("", "*.json")
			if err != nil {
				log.Fatal(err.Error())
				continue
			}

			data, err := json.Marshal(export.Playlist)

			if err != nil {
				log.Fatal(err.Error())
				continue
			}

			if _, err := exportFile.Write(data); err != nil {
				log.Fatal(err.Error())
				continue
			}

			mail := gomail.NewMessage()
			mail.SetHeader("From", "go-open-music@mailtrap.io")
			mail.SetHeader("To", export.TargetEmail)
			mail.SetHeader("Subject", "Export Playlist")
			mail.SetBody("text/html", "You can download this playlist")
			mail.Attach(exportFile.Name())

			d := gomail.NewDialer(config.AppConfig.MAILTRAP_HOST, config.AppConfig.MAILTRAP_PORT, config.AppConfig.MAILTRAP_USERNAME, config.AppConfig.MAILTRAP_PASSWORD)

			if err := d.DialAndSend(mail); err != nil {
				log.Fatal(err)
			}
		}
	}()
}
