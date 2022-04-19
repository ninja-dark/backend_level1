package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	d := net.Dialer{
		Timeout:   time.Second,
		KeepAlive: time.Minute,
	}

	conn, err := d.DialContext(ctx, "tcp", "[::1]:9000")
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(io.Copy(os.Stdout, conn))

	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin) // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
}
