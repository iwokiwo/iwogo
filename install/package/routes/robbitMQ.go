package routes

import (
	"log"

	"github.com/streadway/amqp"
)

func errorWrapper(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SubscribeMessage(ch *amqp.Channel, channelName string) {
	q, err := ch.QueueDeclare(
		channelName, //name
		false,       // durable
		false,       //delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	errorWrapper(err, "Failed to declare a queue")

	msg, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	errorWrapper(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msg {
			//jsonData = []byte(d.Body)
			log.Printf("received as message body: %s", d.Body)

		}

	}()
	log.Printf("waiting for message. to exit press CTRL+C")

	<-forever

}
