package queues

import (
	"errors"

	"github.com/spf13/cobra"
)

type (
	Message interface {
		Data() []byte
		Ack()
	}

	Command *cobra.Command

	Handler func(data []byte) error

	QueueCmd func(handler Handler) Command
)

var (
	ErrCantReadMessage = errors.New("cant read message body")
)
