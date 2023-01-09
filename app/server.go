package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func handleConnection(c net.Conn, wg *sync.WaitGroup) {
	// fmt.Println("Connection established")
	defer wg.Done()
	for {
		buf := make([]byte, 1024)
		len, err := c.Read(buf)

		if err != nil {
			log.Fatal(err)
		}

		str := string(buf[:len])
		// fmt.Println(str)

		switch str {
		case "PING\r\n":
			sendResponse("PONG", c)
		case "QUIT\r\n":
			sendResponse("Goodbye", c)
			c.Close()
		}
		s := fmt.Sprintf("+PONG\r\n")
		c.Write([]byte(s))
	}
}

func sendResponse(res string, conn net.Conn) {
	// time.Sleep(1 * time.Second)
	conn.Write([]byte(res + "\n"))
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Error to listen 6379")
		os.Exit(1)
	}

	defer ln.Close()
	wg := new(sync.WaitGroup)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go handleConnection(conn, wg)
		wg.Wait()
		conn.Close()
	}
}
