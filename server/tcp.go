package server

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
)

func tcpServer(ctx context.Context, lis net.Listener) {
	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down TCP server")
			return
		default:
		}

		conn, err := lis.Accept()
		if err != nil {
			log.Printf("failed to accept: %v\n", err)
			continue
		}

		go func() {
			defer conn.Close()

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					return
				default:
				}

				ln := scanner.Bytes()
				log.Printf("received %d bytes %q from TCP %s\n", len(ln), ln, conn.RemoteAddr())

				if _, err := io.WriteString(conn, "Reply from TCP server\n"); err != nil {
					log.Printf("failed to write to connection: %v\n", err)
				}
			}

			log.Println("connection closed")
		}()
	}
}
