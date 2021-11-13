package main

import (
	"fmt"

	"github.com/fdelbos/queue_handler/queues/amqp"
	"github.com/fdelbos/queue_handler/queues/pubsub"

	"github.com/fdelbos/queue_handler/queues/http"

	"github.com/fdelbos/queue_handler"
)

func main() {
	queue_handler.NewListener().
		Name("basic").
		Description("A very basic example").
		Register(http.Queue()).
		Register(amqp.Queue()).
		Register(pubsub.Queue()).
		Listen(handler)
}

func handler(data []byte) error {
	fmt.Println("Hello from the basic handler!")
	return nil
}
