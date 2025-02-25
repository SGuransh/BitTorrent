package main
import "net"
import "fmt"
import "log"

type Server struct {
	listenAddr string
	ln 		   net.Listener
	quitch    chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch: make(chan struct{}),
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
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Failed to read from connection: ", err)
			continue
		}

		msg := buf[:n]
		fmt.Println("Received message: ", string(msg))
	}
}

func main() {
	server := NewServer(":3000")
	log.Fatal(server.Start())
}