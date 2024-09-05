package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "mydb"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		message TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":5555")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Listening on :5555")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Received message: %s", message)

		sqlStatement := `INSERT INTO messages (message, last_updated) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET message = $1, last_updated = $2`
		_, err = db.Exec(sqlStatement, message, time.Now())
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println("Message inserted/updated successfully!")
	}
}
