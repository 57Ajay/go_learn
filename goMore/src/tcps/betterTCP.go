package tcps

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server_ struct {
	listenAddr string
	ln         net.Listener
	quit       chan struct{}
	wg         *sync.WaitGroup
	timeout    time.Duration
}

func NewServer_(listenAddr string, timeout time.Duration) *Server_ {
	return &Server_{
		listenAddr: listenAddr,
		quit:       make(chan struct{}),
		wg:         &sync.WaitGroup{},
		timeout:    timeout,
	}
}

func (s *Server_) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	log.Printf("Server listening on %s", s.listenAddr)
	defer func() {
		if err := s.ln.Close(); err != nil {
			fmt.Println("Error closing listener: ", err)
		}
		fmt.Println("Listener closed")
	}()
	go s.acceptConnections()

	// These are for gracefull shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigchan:
		log.Println("Shutting down gracefully...")
		close(s.quit) // Signal to stop accepting connections
		s.wg.Wait()   // Wait for all connections to close
		log.Println("Server stopped.")
	case <-s.quit: // In case quit is closed elsewhere
		s.wg.Wait()
		log.Println("Server stopped.")
	}

	return nil
}

func (s *Server_) acceptConnections() {
	defer close(s.quit) // Signal quit when done accepting
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("Accept Error: ", err)
				continue
			}
		}
		s.wg.Add(1) // THis is incrementing the wait group
		go s.handleConnection(conn)
	}
}

func (s *Server_) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection: ", err)
		}
		s.wg.Done() // This is decrementing the wait group
		log.Printf("Connection closed from %s", conn.RemoteAddr())
	}()
	log.Printf("Connection from %s established.", conn.RemoteAddr())

	if s.timeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(s.timeout)); err != nil {
			log.Printf("Set deadline error: %v", err)
			return

		}
	}

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Read Error: ", err)
			}
			return
		}
		log.Printf("Received from %s: %s", conn.RemoteAddr(), string(buf[:n]))
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("Write Error: ", err)
		}
	}
}

func BetterMain() {
	timeOut := 100 * time.Second
	s := NewServer_(":8000", timeOut)
	if err := s.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
