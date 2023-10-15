package main

import (
	"fmt"
	"net"
)

type Client struct {
	conn     *net.TCPConn
	username string
}

var clients []Client

func main() {
	serverAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println(err)
		return
	}

	serverListener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer serverListener.Close()

	fmt.Println("Server is listening on 127.0.0.1:12345")

	for {
		clientConn, err := serverListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle each connection in a separate Goroutine
		go handleConnection(clientConn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	client := Client{conn: conn}

	// Prompt the user for a username
	client.username = getUsername(conn)
	clients = append(clients, client)
	fmt.Printf("%s has joined the chat.\n", client.username)

	for {
		// Read data from the client
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("%s has left the chat.\n", client.username)
			removeClient(client)
			return
		}

		receivedMessage := string(buf[:n])
		fmt.Printf("%s: %s", client.username, receivedMessage)

		// Broadcast the message to all clients
		broadcastMessage(client, receivedMessage)
	}
}

func getUsername(conn *net.TCPConn) string {
	conn.Write([]byte("Enter your username: "))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	return string(buf[:n])
}

func broadcastMessage(sender Client, message string) {
	for _, client := range clients {
		if client.conn != sender.conn {
			client.conn.Write([]byte(sender.username + ": " + message))
		}
	}
}

func removeClient(client Client) {
	for i, c := range clients {
		if c.conn == client.conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}
