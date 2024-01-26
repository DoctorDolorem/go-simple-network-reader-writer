package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var port string

func defineFlags() {
	flag.StringVar(&port, "p", "9000", "port to listen on")
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

	source := ":" + port
	listener, err := net.Listen("tcp", source)
	if err != nil {
		log.Fatal("can't open port", err)
	} else {
		fmt.Println("Listening on port", port)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	} else {
		remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
		fmt.Println("Connection accepted from", remoteAddr.IP)
	}

	handleConnection(conn)
}
