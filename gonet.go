package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var p string

func defineFlags() {
	flag.StringVar(&p, "port", "9000", "port to listen on")
	flag.Parse()
}

func handleConnection(conn net.Conn) {
	_, err := io.Copy(conn, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defineFlags()

	fmt.Println("Listening on port ", p)
	source := ":" + p
	listener, err := net.Listen("tcp", source)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	handleConnection(conn)
}
