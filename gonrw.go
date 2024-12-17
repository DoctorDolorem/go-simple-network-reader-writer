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
	port = ":" + port
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	//var exit bool
	// Goroutine for sending commands
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			cmdStr := scanner.Text() + "\n"
			if cmdStr == "exit\n" {
				_, err := conn.Write([]byte("exit" + "\n"))
				if err != nil {
					log.Println("error writing to connection:", err)
				}
				fmt.Print("Connection closed\n")
				//conn.Close()
				os.Exit(0)

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
	}()

	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		fmt.Printf("error reading from connectionAAA: %s", err)
	}
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

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	fmt.Println("Connection accepted from", remoteAddr.IP)

	handleConnection(conn)

}
