package server

import (
	"log"
	"net"
)

func udpServer(lis net.PacketConn, echoEnabled bool) {
	for {
		buf := make([]byte, 1024)
		n, addr, err := lis.ReadFrom(buf)
		if err != nil {
			return
		}

		go func() {
			log.Printf("received %d bytes %q from UDP %s\n", n, buf[:n], addr.String())
			if !echoEnabled {
				return
			}

			echo := append([]byte("ECHO: "), buf[:n]...)

			if _, err := lis.WriteTo(echo, addr); err != nil {
				log.Printf("failed to write to connection: %v\n", err)
			}
		}()
	}
}
