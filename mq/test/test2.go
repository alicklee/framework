package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, _ := amqp.Dial("amqp://admin:admin@132.232.28.202:5672/")
	defer conn.Close()

	ch, err := conn.Channel()
	fmt.Println(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"gameserver1", // name
		"topic",       // type
		false,         // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	ch.QueueBind("gq1", "s1", "gameserver1", false, nil)
	body := "sssss"
	err = ch.Publish(
		"gameserver1", // exchange
		"s1",          // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
}
