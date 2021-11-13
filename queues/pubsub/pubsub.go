package pubsub

import (
	"fmt"

	"github.com/fdelbos/queue_handler/queues"
	"github.com/spf13/cobra"
)

func command(handler queues.Handler) queues.Command {
	cmd := &cobra.Command{
		Use:   "pubsub",
		Short: "Runs as a Google Cloud Pub/Sub listener",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("pubsub")
		},
	}

	// cmd.Flags().IntVarP(&port, "port", "p", port, "Port to listen on")
	return cmd
}

func Queue() queues.QueueCmd {
	return command
}
