package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	messages = make(chan string)
)

type Server struct {
	listener    net.Listener
	Connections chan net.Conn
}

func NewServer(address string) Server {
	lister, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}

	connChan := make(chan net.Conn)

	return Server{
		listener:    lister,
		Connections: connChan,
	}

}

func (s Server) Start() {
	log.Printf("server started on %s\n", s.listener.Addr())

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		s.Connections <- conn
	}

}

func main() {
	srv := NewServer(":8001")
	go srv.Start()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			log.Println("start graceful")
			wg.Wait()
			log.Println("stop graceful")
			return
		case conn := <-srv.Connections:
			wg.Add(1)
			go handleConn(ctx, conn, &wg)
			go messager()
		}
	}
}

func handleConn(ctx context.Context, c net.Conn, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		c.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			io.WriteString(c, "Bye!")
			return
		case msg := <-messages:
			io.WriteString(c, msg)

		default:
			io.WriteString(c, time.Now().Format("15:04:05\n\r"))
		}
		time.Sleep(1 * time.Second)
	}
}

func messager() {
	for {
		m := bufio.NewReader(os.Stdin)
		r, err := m.ReadString('\n')
		if err != nil {
			return
		}
		messages <- r
	}
}
