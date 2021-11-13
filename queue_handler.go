package queue_handler

import (
	"errors"

	"github.com/fdelbos/queue_handler/queues"
	"github.com/spf13/cobra"
)

type (
	Listener struct {
		name        *string
		description *string
		commands    []queues.QueueCmd
	}
)

var (
	ErrNoCommandRegistered = errors.New("no command registered")
)

func NewListener() Listener {
	return Listener{}
}

func (l Listener) Name(name string) Listener {
	l.name = &name
	return l
}

func (l Listener) Description(desc string) Listener {
	l.description = &desc
	return l
}

func (l Listener) Register(q queues.QueueCmd) Listener {
	if l.commands == nil {
		l.commands = []queues.QueueCmd{}
	}
	l.commands = append(l.commands, q)
	return l
}

func (l Listener) Listen(handler queues.Handler) error {
	name := "unknown"
	desc := ""

	if l.name != nil {
		name = *l.name
	}
	if l.description != nil {
		desc = *l.description
	}

	rootCmd := &cobra.Command{
		Use:   name,
		Short: desc,
	}

	if l.commands == nil || len(l.commands) == 0 {
		return ErrNoCommandRegistered
	}
	for _, cmd := range l.commands {
		rootCmd.AddCommand(cmd(handler))
	}

	return rootCmd.Execute()
}
