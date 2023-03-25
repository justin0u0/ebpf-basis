package tcpserver

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var remoteAddr string

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tcpserver",
		Run: run,
	}

	cmd.Flags().StringVarP(&remoteAddr, "remote", "r", remoteAddr, "remote address to connect to")

	return cmd
}

func run(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()

	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Fatalf("failed to dial: %v\n", err)
	}
	log.Println("local address:", conn.LocalAddr())

	go client(conn)
	go server(conn)

	<-ctx.Done()

	log.Println("shutting down...")

	if err := conn.Close(); err != nil {
		log.Printf("failed to close connection: %v\n", err)
	}
}

func client(conn net.Conn) {
	// client wait for os.Stdin and send it to the remote server
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		ln := s.Bytes()
		log.Printf("sending %d bytes %q to TCP %s\n", len(ln), ln, conn.RemoteAddr())
		ln = append(ln, '\n')

		if _, err := conn.Write(ln); err != nil {
			log.Printf("failed to write to connection: %v\n", err)
			continue
		}
	}
}

func server(conn net.Conn) {
	// server listen the echo message from the remote server
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Bytes()
		log.Printf("received %d bytes %q from TCP %s\n", len(ln), ln, conn.RemoteAddr())
	}
}
