package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

const LISTEN_ADDR = "localhost:4000"

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socket_handler))
	err := http.ListenAndServe(LISTEN_ADDR, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type socket struct {
	io.ReadWriter
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Patience is a virtue...")
}

func socket_handler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool)}
	go match(s)
	<-s.done
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
	fmt.Fprint(a, "Welcome, pilgrim.\n")
	fmt.Fprint(b, "Welcome, pilgrim.\n")

	errc := make(chan error, 1)

	go cp(a, b, errc)
	go cp(b, a, errc)

	if err := <-errc; err != nil {
		log.Fatal(err)
	}
	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}
