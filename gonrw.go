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
var exit = false

func defineFlags() {
	flag.StringVar(&port, "p", "9000", "port to listen on")
	flag.Parse()
	port = ":" + port

}

func handleConnection(conn net.Conn) error {
	// Goroutine for sending commands
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			cmdStr := scanner.Text() + "\n"
			if cmdStr == "close;\n" {
				fmt.Println("Closing connection")
				exit = true
				conn.Write([]byte("exit\n"))
				break
			}
			_, err := conn.Write([]byte(cmdStr))
			if err != nil {
				log.Println("error writing to connection:", err)
			}
		}
		err := scanner.Err()
		if err != nil {
			log.Println("error reading from stdin:", err)
		}

		if exit {
			return
		}
	}()

	if exit {
		return nil
	}

	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		return fmt.Errorf("error reading from connection: %s", err)
	}
	return nil
}

func main() {
	defineFlags()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("can't start listener: ", err)
	}

	fmt.Println("Listening on port", port)
	fmt.Println("type 'close;' at any moment to close the connection")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("error accepting incoming connection: %s", err)
		conn.Close()
	}
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	fmt.Println("Connection accepted from", remoteAddr.IP)

	err = handleConnection(conn)
	if err != nil {
		fmt.Printf("error while handling connection: %s\n", err)
		conn.Close()
	}
	if exit {
		conn.Close()
		fmt.Printf("Connection closed by user\n")
		os.Exit(0)
	}
}
