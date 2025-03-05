package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	fmt.Println("Listening on: ", s.listenAddr)
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err)
			continue
		}

		fmt.Println("Accepted connection from: ", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	var message []byte // To accumulate characters

	for {
		buf := make([]byte, 1) // Read one byte at a time
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				// EOF error is expected when the client disconnects
				fmt.Println("Connection closed by client: ", conn.RemoteAddr())
				break
			}
			fmt.Println("Failed to read from connection: ", err)
			continue
		}

		// Append the character to the message buffer
		message = append(message, buf[:n]...)

		// Check for CRLF (\r\n) to signal the end of the message
		if len(message) >= 2 && message[len(message)-2] == '\r' && message[len(message)-1] == '\n' {
			// Remove the \r\n from the message to process the actual message
			message = message[:len(message)-2]

			// Send the full message to the msgch channel
			s.msgch <- Message{
				from:    conn.RemoteAddr().String(),
				payload: message,
			}

			// Respond back to the client immediately
			conn.Write([]byte("Message received\r\n"))

			// Reset the message buffer for the next input
			message = nil
		}
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgch {
			fmt.Printf("Message from %s: %s\n", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}
