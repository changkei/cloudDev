// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

type client struct {
	chann chan<- string // an outgoing message channel
	name  string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients

	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.chann <- msg
			}

		case cli := <-entering:
			clients[cli] = true

			clientList := "Present Clients: "

			if len(clients) == 1 {
				clientList = clientList + "You are currently alone\n"
			} else if len(clients) > 1 {
				clientList = clientList + "There are currently " +
					strconv.Itoa(len(clients)-1) + " other clients\n"
			}

			for c := range clients {
				clientList = clientList + c.name + "\n"
			}

			cli.chann <- clientList

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.chann)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	newCli := client{chann: ch}
	go clientWriter(conn, ch)

	ch <- "Who are you?"

	//newCli.name = conn.RemoteAddr().String()
	oinput := bufio.NewScanner(conn)
	oinput.Scan()
	newCli.name = oinput.Text()

	ch <- "You are " + newCli.name + "\n"

	messages <- newCli.name + " has arrived\n"
	entering <- newCli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- newCli.name + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- newCli
	messages <- newCli.name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
