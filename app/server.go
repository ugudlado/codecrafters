package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer connection.Close()

	buffer := make([]byte, 1024)
	_, err = connection.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read data from connection")
		os.Exit(1)
	}

	data := string(buffer)

	streams := strings.Split(data, "\r\n")

	actualUrl := strings.Split(streams[0], " ")[1]

	expectedUrl := "/"

	if strings.Compare(actualUrl, expectedUrl) == 0 {
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
