package server

import (
	"log"
	"net"
	"sync"

	"github.com/spf13/cobra"
)

var (
	tcpPort     = ":8080"
	udpPort     = ":8081"
	echoEnabled = false
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
		Run: run,
	}

	cmd.Flags().StringVarP(&tcpPort, "tcp", "t", tcpPort, "TCP port to listen on")
	cmd.Flags().StringVarP(&udpPort, "udp", "u", udpPort, "UDP port to listen on")
	cmd.Flags().BoolVarP(&echoEnabled, "echo", "e", echoEnabled, "enable echo server")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	tcpLis, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalf("failed to listen to TCP on port %s: %v", tcpPort, err)
	}
	log.Println("listening TCP on", tcpPort)

	udpLis, err := net.ListenPacket("udp", udpPort)
	if err != nil {
		log.Fatalf("failed to listen to UDP on port %s: %v", udpPort, err)
	}
	log.Println("listening UDP on", udpPort)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		tcpServer(tcpLis, echoEnabled)
		wg.Done()
	}()
	go func() {
		udpServer(udpLis, echoEnabled)
		wg.Done()
	}()

	<-ctx.Done()

	log.Println("shutting down...")

	if err := tcpLis.Close(); err != nil {
		log.Printf("failed to close TCP listener: %v\n", err)
	}
	if err := udpLis.Close(); err != nil {
		log.Printf("failed to close UDP listener: %v\n", err)
	}

	wg.Wait()
}
