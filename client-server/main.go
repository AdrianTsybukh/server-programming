package main

import (
	"log"
	"net"
	"bufio"
	"fmt"
	"strings"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Printf("Connection %s accepted.\n", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	requestLine = strings.TrimSpace(requestLine)
	fmt.Printf("Raw request line: %s\n", requestLine)

	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		log.Println("Invalid HTTP request.")
		return
	}

	method := parts[0]
	path := parts[1]
	httpVersion := parts[2]

	fmt.Printf("PARSED Method: %s. Path: %s. HTTP version: %s.\n", method, path, httpVersion)

	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		header = strings.TrimSpace(header)
		if header == "" {
			break
		}

		fmt.Printf("Header: %s\n", header)
	}

	var response string
	if method == "GET" {
		switch path {
		case "/hello":
			body := "Hello, world!"
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", len(body), body)
		case "/goodbye":
			body := "Goodbye, world.."
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", len(body), body)
		default:
			body := "Not Found"
			response = fmt.Sprintf("HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", len(body), body)
		}
	} else {
		body := "Method Not Allowed"
		response = fmt.Sprintf("HTTP/1.1 405 Method Not Allowed\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s", len(body), body)
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println(err)
	}
}


func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	log.Println("Listening to port 8080.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		
		go handleConn(conn)
	}
}

