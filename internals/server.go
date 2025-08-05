package internals

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var connections []net.Conn

func handleClient(conn net.Conn) {
	connections = append(connections, conn)
	reader := bufio.NewReader(conn)

	username, err := reader.ReadString('\n')
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}
	username = strings.TrimSpace(username)

	for i, c := range connections {
		_, err := c.Write([]byte(fmt.Sprintf("[!] %s has joined!\n", string(username))))
		if err != nil {
			c.Close()
			connections = append(connections[i:], connections[i+1:]...)
			i--
		}
	}

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("%s closed", conn.RemoteAddr().String())
			break
		}

		for i, c := range connections {
			if c != conn {
				_, err := c.Write([]byte(fmt.Sprintf("%s: %s", username, message)))
				if err != nil {
					c.Close()
					connections = append(connections[i:], connections[i+1:]...)
					i--
				}
			}
		}
	}
}

func CreateServer(port int) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("[Server IP]$: ")
	addr := ""
	if scanner.Scan() {
		addr = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if addr == "" {
		fmt.Println("?")
		return
	}

	host := net.JoinHostPort(addr, strconv.Itoa(port))
	l, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	log.Println("Server listening on", host)

	for {
		conn, err := l.Accept()
		if err != nil {
			if len(connections) <= 0 {
				continue
			}

			for _, c := range connections {
				c.Write([]byte(err.Error() + "\n"))
			}
			continue
		}

		log.Printf("%s connected", conn.RemoteAddr().String())

		go handleClient(conn)
	}
}

func Shutdown() {
	for _, c := range connections {
		c.Close()
	}
}
