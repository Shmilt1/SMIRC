package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"smirc/internals"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	isServer := flag.Bool("server", false, "Start a custom SMIRC server")

	flag.Parse()

	if *isServer {
		go internals.CreateServer(8080)
	} else {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("[Server Host]$: ")
		host := ""
		if scanner.Scan() {
			host = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}

		if host == "" {
			fmt.Println("?")
			return
		}

		go internals.ConnectClient(host)
	}

	<-quit
	log.Println("Shutting down...")

	if *isServer {
		internals.Shutdown()
	}

	os.Exit(0)
}
