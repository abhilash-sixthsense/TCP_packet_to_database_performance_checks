package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5555")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	messageCount := 0
	for {
		messageCount++
		timestamp := time.Now().Format(time.RFC3339)
		message := fmt.Sprintf("Message %d sent at %s\n", messageCount, timestamp)
		_, err := fmt.Fprintf(conn, message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Sent message:", message)
		// time.Sleep(1 * time.Second)
	}
}
