package config

import (
	"fmt"

	"github.com/streadway/amqp"
)

func NewQueue() *amqp.Connection {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/go-open-music", AppConfig.RABBITMQ_USERNAME, AppConfig.RABBITMQ_PASSWORD, AppConfig.RABBITMQ_HOST, AppConfig.RABBITMQ_PORT))
	if err != nil {
		panic(err.Error())
	}

	return conn
}
