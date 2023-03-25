package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/justin0u0/ebpf-basis/bpfgo"
	"github.com/justin0u0/ebpf-basis/client"
	"github.com/justin0u0/ebpf-basis/server"
	"github.com/justin0u0/ebpf-basis/server/tcpserver"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{}

	cmd.AddCommand(
		server.Command(),
		tcpserver.Command(),
		bpfgo.Command(),
		client.Command(),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := cmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
