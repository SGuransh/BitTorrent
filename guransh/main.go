package main
import "net"

type Server struct {
	listenAddr string
	ln 		   net.Listener
	quitch    chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr
	}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	<-s.quitch

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err)
			continue
		}
	}
}

func (s *Server) readLoop(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Failed to read from connection: ", err)
			return
		}
	}
}

func main() {
	println("Hello, World!")
}