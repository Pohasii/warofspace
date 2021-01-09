package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type message struct {
	addr *net.UDPAddr
	text []byte
}

type Messages chan message

type ValidUser map[string]time.Time // *net.UDPAddr

var secret string = "qwerty"

type network struct {
	secret string
	maxMessageSize int
	timeOutConnections time.Duration
	serverAddress *net.UDPAddr
	connections map[string]time.Time
	socket *net.UDPConn
	status bool
}

func (the *network) reader(MessagesFromUser Messages) {
	defer the.socket.Close()

	var sms = make([]byte, the.maxMessageSize)

	for the.status {

		size, caddr, err := the.socket.ReadFromUDP(sms)
		if err != nil {
			log.Println(err)
		}

		if the.validation(caddr, &sms, &size); size > 0 {
			MessagesFromUser <- message{caddr, sms[:size]}
		}
	}
}

func (the *network) validation(userAddr *net.UDPAddr, letter *[]byte, size *int) {

	// check user in verified Users or validation secret
	if _, ok := the.connections[userAddr.String()]; ok || secret == string((*letter)[:6]) {

		// update time last message
		the.connections[userAddr.String()] = time.Now()

		// del secrets byte from message
		*letter = (*letter)[6:]
		if string(*letter) == "ping" {
			*size = 0
		}
	} else {
		*size = 0
	}
}

//
func (the *network) sender(MessagesToUser Messages) {
	for letter := range MessagesToUser {
		_, err := the.socket.WriteTo(letter.text, letter.addr)
		if err != nil {
			log.Println(err)
		}
	}
}

func (the *network) handler(MessagesFromUser, MessagesToUser Messages) {

	// slice for offline users
	offlineUsers := make([]string, 0, 300)

	for mess := range MessagesFromUser {

		for client, date := range the.connections {

			elapsed := time.Now().Sub(date)

			if client != mess.addr.String() && elapsed < the.timeOutConnections { //
				addr, err := net.ResolveUDPAddr("udp", client)
				if err != nil {
					log.Println(err)
				}
				MessagesToUser <- message{addr, mess.text}
			}

			if elapsed > the.timeOutConnections {
				offlineUsers = append(offlineUsers, client)
			}
		}

		// remove offline users
		for _, disconnect := range offlineUsers {
			delete(the.connections, disconnect)
		}
		offlineUsers = nil // make([]string, 0, 1000)
		offlineUsers = make([]string, 0, 300)
	}
}

func (the *network) netStart() {
	fmt.Println("Server Start")

	// init server
	// sAddr := "0.0.0.0:55442" //localhost:0 192.168.0.52  // "0.0.0.0:55442" "192.168.0.52:55442"

	go handler(Users, MessagesFromUsers, MessagesToUsers)
	go sender(Conn, MessagesToUsers)
	reader(Conn, MessagesFromUsers, Users)
}

func (the *network) SetServerAddr (addr string) {

	adr, err := net.ResolveUDPAddr("udp", addr) //192.168.0.52:12345
	if err != nil {
		log.Println(err)
	}

	the.serverAddress = adr
}

func (the *network) OpenSocket () {

	socket, err := net.ListenUDP("udp", the.serverAddress)
	if err != nil {
		log.Println(err)
	}

	the.socket = socket
}