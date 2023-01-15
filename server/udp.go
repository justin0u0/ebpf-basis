package server

import (
	"context"
	"log"
	"net"
)

func udpServer(ctx context.Context, lis net.PacketConn) {
	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down UDP server")
			return
		default:
		}

		buf := make([]byte, 1024)
		n, addr, err := lis.ReadFrom(buf)
		if err != nil {
			log.Printf("failed to read: %v\n", err)
			continue
		}

		go func() {
			log.Printf("received %d bytes %q from UDP %s\n", n, buf[:n], addr.String())

			if _, err := lis.WriteTo([]byte("Reply from UDP server\n"), addr); err != nil {
				log.Printf("failed to write to connection: %v\n", err)
				return
			}
		}()
	}
}
