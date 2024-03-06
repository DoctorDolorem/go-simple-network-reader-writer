package main

import (
	"bufio"
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
	defer conn.Close()

	// Goroutine for sending commands
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			cmdStr := scanner.Text() + "\n"
			_, err := conn.Write([]byte(cmdStr))
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	// Goroutine for receiving output
	_, err := io.Copy(os.Stdout, conn)
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
