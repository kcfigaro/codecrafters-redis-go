package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func handleBufferConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	var msg []string

	for {
		// read buf bytes from conn
		n, err := conn.Read(buf)
		// fmt.Printf("%d byte, data: %s\n", n, buf[:n])
		msg = append(msg, string(buf[:n]))
		if err == io.EOF {
			fmt.Println("end")
			break
		} else if err != nil {
			fmt.Println("error", err)
			break
		}

		var fields []string
		fields = strings.Fields(msg[0])
		// fmt.Println("fields: ", fields)
		switch fields[2] {
		case "ping":
			responseConnection("PONG", conn)
		case "echo":
			responseConnection(fields[4], conn)
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
		go handleBufferConn(conn)
	}
}
