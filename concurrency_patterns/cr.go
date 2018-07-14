package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const LISTEN_ADDR = "localhost:4000"

func main() {
	l, err := net.Listen("tcp", LISTEN_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go match(c)
	}
}

var p = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Patience is a virtue...\n> ")
	select {
	case p <- c:
	case pm := <-p:
		chat(pm, c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintf(a, "Welcome, pilgrim.\n")
	fmt.Fprintf(b, "Welcome, pilgrim.\n")

	errc := make(chan error, 1)

	go cp(a, b, errc)
	go cp(b, a, errc)

	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}
