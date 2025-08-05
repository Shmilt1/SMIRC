package internals

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var connections []net.Conn

func handleClient(conn net.Conn) {
	conn.SetDeadline(time.Now().Add(time.Second * 3))

	connections = append(connections, conn)

	username := make([]byte, 256)
	n, err := conn.Read(username)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	if n <= 0 {
		conn.Write([]byte("[*] Please set username!"))
		return
	}

	for _, c := range connections {
		c.Write([]byte(fmt.Sprintf("[*] %s has joined!", string(username))))
	}

	message := make([]byte, 256)
	total := 0
	for {
		n, err := conn.Read(message[total:])
		if err != nil {
			log.Println(err)
			continue
		}

		if n <= 0 {
			continue
		}

		total += n
	}
}

func CreateServer(port int) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("[Server IP]$: ")
	host := ""
	if scanner.Scan() {
		host = scanner.Text()
	}
	log.Println("Chosen server ip", host)

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if host == "" {
		fmt.Println("?")
		os.Exit(0)
	}

	l, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			if len(connections) <= 0 {
				continue
			}

			for _, c := range connections {
				c.Write([]byte(err.Error()))
			}
			continue
		}

		go handleClient(conn)
	}
}
