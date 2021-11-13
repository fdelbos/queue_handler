package http

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	shttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/fdelbos/queue_handler/queues"
	"github.com/spf13/cobra"
)

type (
	httpQueue struct {
		port    int
		handler queues.Handler
	}
)

const DefaultPort = 80

func command(handler queues.Handler) queues.Command {
	q := httpQueue{
		port:    DefaultPort,
		handler: handler,
	}

	cmd := &cobra.Command{
		Use:   "http",
		Short: "Runs as an http service",
		Run: func(cmd *cobra.Command, args []string) {
			q.run()
		},
	}

	cmd.Flags().IntVarP(&q.port, "port", "p", DefaultPort, "Port to listen on")
	return cmd
}

func Queue() queues.QueueCmd {
	return command
}

func (q httpQueue) run() {
	ctx, cancel := context.WithCancel(context.Background())
	addr := fmt.Sprintf(":%d", q.port)

	httpServer := shttp.Server{
		Addr:    addr,
		Handler: q.handlerFunc(),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http server error")
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	<-signalChan
	log.Print("os.Interrupt, shutting down...\n")

	go func() {
		<-signalChan
		log.Fatal().Msg("os.Kill - terminating...\n")
	}()

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(gracefullCtx); err != nil {
		log.Error().Err(err).Msg("shutdown error")
		defer os.Exit(1)
		return
	} else {
		log.Info().Msg("gracefully stopped\n")
	}

	cancel()

	defer os.Exit(0)
}

func setStatusError(w shttp.ResponseWriter) {
	w.WriteHeader(shttp.StatusInternalServerError)
}

func setStatusOK(w shttp.ResponseWriter) {
	w.WriteHeader(shttp.StatusOK)
}

func (q httpQueue) handlerFunc() shttp.HandlerFunc {
	return func(w shttp.ResponseWriter, r *shttp.Request) {
		if data, err := io.ReadAll(r.Body); err != nil {
			log.Error().Err(err).Msg("cant read request body")

		} else if err := q.handler(data); err != nil {
			log.Error().Err(err).Msg("request returned an error")

		} else {
			setStatusOK(w)
			return
		}
		setStatusError(w)
	}
}
