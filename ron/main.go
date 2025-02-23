package main

import "net"

type Server struct {
	listenAddr string
	ln 		   net.Listener
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server {
		listenAddr: listenAddr,
		quitch: 	make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tpc", s.listenAddr)
	"""
	net.Listen() creates a network listener (AKA server) and you
	can respond to incoming requests with acceptLoop below.

	How it works:
	1. OS binds given address and port
	2. Program listens for incoming requests
	3. Once a client connects, you can accept using ln.Accept()
	"""
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	<-s.quitch  // Blocks here until quitch is closed
	"""
	Understanding <-s.quitch
	s.quitch is a channel of type chan struct{}.
	<-s.quitch means 'receive from the channel'.
	Since quitch is an unbuffered channel, the operation will block until:
	1. A value is sent to s.quitch → s.quitch <- struct{}{}.
	2. The channel is closed → close(s.quitch).
	"""

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Acccept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go s.readLoop(conn)  // go keyword used to start thread routine - Multithreading
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

		msg := buf[:n]
		fmt.Println(string(msg))
	}
}

func main() {
}