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

	client.username = getUsername(conn)
	clients = append(clients, client)
	notify(client, "has joined the chat.\n")
	fmt.Printf("%s has joined the chat.\n", client.username)

	for {
		// Read data from the client
		n, err := conn.Read(buf)
		if err != nil {
			notify(client, "has left the chat.\n")

			removeClient(client)
			return
		}

		receivedMessage := string(buf[:n])
		fmt.Printf("%s: %s", client.username, receivedMessage)

		// Broadcast
		broadcastMessage(client, receivedMessage)
	}
}

func notify(sender Client, message string) {
	for _, client := range clients {
		if client.conn != sender.conn {
			client.conn.Write([]byte(sender.username + " " + message))
			return
		}
	}
}

func getUsername(conn *net.TCPConn) string {
	// conn.Write([]byte("Enter your username: "))
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
