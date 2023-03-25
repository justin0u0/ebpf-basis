package server

import (
	"bufio"
	"log"
	"net"
)

func tcpServer(lis net.Listener, echoEnabled bool) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				ln := scanner.Bytes()
				log.Printf("received %d bytes %q from TCP %s\n", len(ln), ln, conn.RemoteAddr())

				if !echoEnabled {
					continue
				}

				ln = append(ln, '\n')
				echo := append([]byte("ECHO: "), ln...)

				if _, err := conn.Write(echo); err != nil {
					log.Printf("failed to write to connection: %v\n", err)
					continue
				}
			}

			log.Println("connection closed")
		}()
	}
}
