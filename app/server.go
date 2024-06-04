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

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	for {
		handleRequest(connection)
	}
}

func handleRequest(connection net.Conn) {
	buffer := make([]byte, 1024)
	_, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read data from connection")
		os.Exit(1)
	}

	data := string(buffer)

	streams := strings.Split(data, "\r\n")

	actualUrl := strings.Split(streams[0], " ")[1]

	responseMessage := "HTTP/1.1 404 Not Found\r\n\r\n"

	if strings.HasPrefix(actualUrl, "/echo/") {
		echoMessage := strings.Replace(actualUrl, "/echo/", "", 1)
		responseMessage = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoMessage), echoMessage)
	} else if strings.HasPrefix(actualUrl, "/user-agent") {
		for i := 1; i < len(streams); i++ {
			if strings.HasPrefix(streams[i], "User-Agent") {
				header := strings.Replace(streams[i], "User-Agent: ", "", 1)
				responseMessage = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(header), header)
				break
			}
		}
	} else if strings.Compare(actualUrl, "/") == 0 {
		responseMessage = "HTTP/1.1 200 OK\r\n\r\n"
	}
	print(actualUrl)
	connection.Write([]byte(responseMessage))
}
