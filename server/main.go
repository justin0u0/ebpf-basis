package server

import (
	"log"
	"net"

	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	tcpLis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen to TCP on port 8080: %v", err)
	}
	defer tcpLis.Close()

	log.Println("listening TCP on :8080")

	udpLis, err := net.ListenPacket("udp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen to UDP on port 8081: %v", err)
	}
	defer udpLis.Close()

	log.Println("listening UDP on :8081")

	ctx := cmd.Context()

	go tcpServer(ctx, tcpLis)
	go udpServer(ctx, udpLis)

	<-ctx.Done()
}
