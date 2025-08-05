package internals

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func ConnectClient(host string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("[Username]$: ")
	username := ""
	if scanner.Scan() {
		username = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if username == "" {
		fmt.Println("?")
		return
	}

	_, err = conn.Write([]byte(username + "\n"))
	if err != nil {
		log.Fatalln(err)
	}

	reader := bufio.NewReader(conn)
	joined, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(joined)

	go func() {
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(message)
		}
	}()

	for {
		fmt.Printf("[%s]$: ", username)

		if scanner.Scan() {
			input := scanner.Text()

			_, err := conn.Write([]byte(input + "\n"))
			if err != nil {
				log.Fatalln(err)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}
	}
}
