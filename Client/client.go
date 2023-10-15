package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println(err)
		return
	}

	clientConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer clientConn.Close()

	// Get the username from the user
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	username = username[:len(username)-1]
	clientConn.Write([]byte(username))

	// Start a go routine to read and print messages from the server
	go func() {
		// Continuously read and print messages from the server
		buf := make([]byte, 1024)
		for {
			n, err := clientConn.Read(buf)
			if err != nil {
				fmt.Println("Connection to the server is closed.")
				os.Exit(0)
			}

			serverMessage := string(buf[:n])
			fmt.Print(serverMessage)
		}
	}()

	for {
		// Send messages to the server
		message, _ := reader.ReadString('\n')
		clientConn.Write([]byte(message))
	}
}
