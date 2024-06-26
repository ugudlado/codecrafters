package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type ServerOptions struct {
	Directory string
}

func parseArgs(args []string) ServerOptions {
	options := ServerOptions{}
	for i := 0; i < len(args); i++ {
		if args[i] == "--directory" && i+1 < len(args) {
			options.Directory = args[i+1]
		}
	}
	return options
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Read command-line arguments
	args := os.Args[1:] // Exclude the first argument, which is the program name
	if len(args) > 1 {
		if strings.Compare(args[0], "--directory") == 0 {
			if len(args) < 2 {
				fmt.Println("No directory provided")
				os.Exit(1)
			}

			directory := args[1]

			if _, err := os.Stat(directory); os.IsNotExist(err) {
				fmt.Println("Directory does not exist")
				os.Exit(1)
			}
		}
	}

	serverOptions := parseArgs(args)

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

		go handleConnection(connection, serverOptions)
	}
}

func handleConnection(connection net.Conn, serverOptions ServerOptions) {

	buffer := make([]byte, 1024)
	_, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read data from connection")
		os.Exit(1)
	}

	// Remove null bytes from buffer
	data := strings.ReplaceAll(string(buffer), "\x00", "")

	responseMessage := handleRequest(data, serverOptions)
	connection.Write([]byte(responseMessage))
}

func handleRequest(data string, serverOptions ServerOptions) string {
	streams := strings.Split(data, "\r\n")

	actualUrl := strings.Split(streams[0], " ")[1]

	responseMessage := "HTTP/1.1 404 Not Found\r\n\r\n"

	if strings.HasPrefix(streams[0], "POST") {
		if strings.HasPrefix(actualUrl, "/files/") {
			filePath := strings.Replace(actualUrl, "/files/", "", 1)
			filePath = serverOptions.Directory + filePath
			fileContent := streams[len(streams)-1]
			file, err := os.Create(filePath)
			if err != nil {
				responseMessage = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
			}
			file.Write([]byte(fileContent))
			file.Close()
			responseMessage = "HTTP/1.1 201 Created\r\n\r\n"
		}
	} else {
		if strings.HasPrefix(actualUrl, "/files/") {
			filePath := strings.Replace(actualUrl, "/files/", "", 1)
			filePath = serverOptions.Directory + filePath
			if file, err := os.Open(filePath); err == nil {
				fileInfo, _ := file.Stat()
				fileSize := fileInfo.Size()
				fileContent := make([]byte, fileSize)
				file.Read(fileContent)
				responseMessage = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", fileSize, fileContent)
			}
		} else if strings.HasPrefix(actualUrl, "/echo/") {
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
	}
	return responseMessage
}
