package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}

	l, err := cfg.Listen(ctx, "tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	log.Println("im started!")
	go broadcaster()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println(err)
				return
			} else {
				wg.Add(1)
				go handleConn(conn, wg)
				go writer(conn, ctx)
			}
		}
	}()

	<-ctx.Done()

	log.Println("done")
	l.Close()
	wg.Wait()
	log.Println("exit")
}

func handleConn(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	ch := make(chan string)
	entering <- ch
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
	leaving <- ch
	conn.Close()
}

func writer(conn net.Conn, ctx context.Context) {
	ch := make(chan string)
	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			ch <- input.Text()
		}
	}()
	ticker := time.NewTicker(2 * time.Second)

	var message string
	for {
		<-ticker.C
		select {
		case <-ctx.Done():
			return
		case t := <-ch:
			message = t
		default:
			message = time.Now().Format("15:04:05")
		}
		messages <- message
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)

		}
	}

}
