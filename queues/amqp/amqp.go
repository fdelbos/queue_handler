package amqp

import (
	"fmt"

	"github.com/fdelbos/queue_handler/queues"
	"github.com/spf13/cobra"
)

func command(handler queues.Handler) queues.Command {
	cmd := &cobra.Command{
		Use:   "amqp",
		Short: "Runs as an AMQP (RabbitMQ) listener",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("amqp")
		},
	}

	return cmd
}

func Queue() queues.QueueCmd {
	return command
}
