package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func handleConnection(c net.Conn) {
	// fmt.Println("Connection established")
	defer c.Close()
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		str := scanner.Text()
		// fmt.Println(str)

		switch str {
		case "ping":
			responseConnection("PONG", c)
		case "quit":
			responseConnection("QUIT", c)
		}
	}
}

func responseConnection(s string, c net.Conn) {
	c.Write([]byte("+" + s + "\r\n"))
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Error to listen 6379")
		os.Exit(1)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConnection(conn)
	}
}
