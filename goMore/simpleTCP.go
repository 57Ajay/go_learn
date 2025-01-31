package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quit       chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quit:       make(chan struct{}),
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

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan

	fmt.Println("Shutting down server...")
	close(s.quit)
	return nil
}

func (s *Server) acceptLoop() {
	defer s.ln.Close()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.quit:
				fmt.Println("Server shutting down, exiting accept loop.")
				return
			default:
				fmt.Println("accept error: ", err)
				continue
			}
		}
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read Error: ", err)
			return
		}
		msg := buf[:n]
		fmt.Println(string(msg))
	}
}

func main() {
	server := NewServer(":3000")
	if err := server.Start(); err != nil {
		fmt.Println("Error starting server:", err)
	}
	fmt.Println("Server exited.")
}
