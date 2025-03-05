/*
"go run main.go" to run the server

On the other terminal, run "nc localhost 3000" to connect to the server.
You can type messages and hit enter to send them to the server.
*/

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
		/*
			Channels allow goroutines to communicate with each other.
			Here, we create a channel of type []byte with a buffer size of 10.
		*/
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	/*
		net.Listen() creates a network listener (AKA server) and you
		can respond to incoming requests with acceptLoop below.

		How it works:
		1. OS binds given address and port
		2. Program listens for incoming requests
		3. Once a client connects, you can accept using ln.Accept()
	*/
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop() // go keyword used to start thread routine

	<-s.quitch // Blocks here until quitch is closed
	close(s.msgch)

	/*
		Understanding <-s.quitch
		s.quitch is a channel of type chan struct{}.
		<-s.quitch means 'receive from the channel'.
		Since quitch is an unbuffered channel, the operation will block until:
		1. A value is sent to s.quitch → s.quitch <- struct{}{}.
		2. The channel is closed → close(s.quitch).
	*/

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Println("new connection:", conn.RemoteAddr())

		go s.readLoop(conn) // go keyword used to start thread routine - Multithreading
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		conn.Write([]byte("thank you for your message!"))
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgch {
			fmt.Printf("received message from connection (%s):%s\n", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}
