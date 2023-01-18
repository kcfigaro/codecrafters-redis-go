package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func handleBufferConn(kv map[string]string, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	var msg []string

	for {
		// read buf bytes from conn
		n, err := conn.Read(buf)
		// append buf to list of string
		msg = append(msg, string(buf[:n]))
		if err == io.EOF {
			fmt.Println("end")
			break
		} else if err != nil {
			fmt.Println("error", err)
			break
		}

		var fields []string
		// slicing substrings
		fields = strings.Fields(msg[0])
		// fmt.Println("fields: ", fields)

		switch fields[2] {
		case "ping":
			responseConnection("PONG", conn)
		case "echo":
			responseConnection(fields[4], conn)
		case "set":
			kv = setValue(kv, fields[4], fields[6])
			responseConnection("OK", conn)
			if len(fields) > 8 && fields[8] == "px" {
				go expiryVaule(kv, fields[4], fields[10])
			}
		case "get":
			// conn.Write([]byte("+" + kv[fields[4]] + "\r\n"))
			val := getValue(kv, fields[4])
			responseConnection(val, conn)
		}
		// clear the request strings
		msg = make([]string, 0)
	}
}

func expiryVaule(kv map[string]string, key string, expiryTime string) {

	seconds, _ := strconv.Atoi(expiryTime)
	time.Sleep(time.Duration(seconds) * time.Microsecond)
	// fmt.Println("expired ", key)
	delete(kv, key)
}

func setValue(kv map[string]string, key, value string) map[string]string {
	// fmt.Printf("SET: %s, %s", key, value)
	kv[key] = value
	return kv
}

func getValue(kv map[string]string, key string) string {
	// fmt.Println("GET: ", key)
	return kv[key]
}

func responseConnection(s string, c net.Conn) {
	if len(s) == 0 {
		c.Write([]byte("$-1\r\n"))
	}
	c.Write([]byte("+" + s + "\r\n"))
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Error to listen 6379")
		os.Exit(1)
	}

	defer ln.Close()
	kv := make(map[string]string)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleBufferConn(kv, conn)
	}
}
