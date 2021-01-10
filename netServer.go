package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type message struct {
	addr *net.UDPAddr
	text []byte
}

type Messages chan message

type connection map[int]client
// var secret string = "qwerty"
type client struct {
	ID int
	lastActivity time.Time
	addr string
}

type Network struct {
	nextID _ID
	secret string
	maxMessageSize int
	timeOutConnections time.Duration
	serverAddress *net.UDPAddr
	connections connection
	socket *net.UDPConn
	status bool
}

func (the *Network) init () *Network{
	return &Network{
		connections: make(connection, 0),
	}
}

func (the *Network) reader(MessagesFromUser Messages) {
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

func (the *Network) validation(userAddr *net.UDPAddr, letter *[]byte, size *int) {

	// check user in verified Users or validation secret
	if _, ok := the.connections[userAddr.String()]; ok || the.secret == string((*letter)[:6]) {

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
func (the *Network) sender(MessagesToUser Messages) {
	for letter := range MessagesToUser {
		_, err := the.socket.WriteTo(letter.text, letter.addr)
		if err != nil {
			log.Println(err)
		}
	}
}

func (the *Network) handler(MessagesFromUser, MessagesToUser Messages) {

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

func (the *Network) netStart(serverAddress, secret string, maxMessageSize, timeOutConnections int) {
	fmt.Println("Server Start")

	the.setServerAddr(serverAddress)
	the.secret = secret
	the.setMessageSize(maxMessageSize)
	the.setTimeOutConn(timeOutConnections)
	the.status = true
	the.OpenSocket()

	// init server
	// sAddr := "0.0.0.0:55442" //localhost:0 192.168.0.52  // "0.0.0.0:55442" "192.168.0.52:55442"

	//go handler(Users, MessagesFromUsers, MessagesToUsers)
	//go sender(Conn, MessagesToUsers)
	//reader(Conn, MessagesFromUsers, Users)
}

func (the *Network) setServerAddr(addr string) {

	adr, err := net.ResolveUDPAddr("udp", addr) //192.168.0.52:12345
	if err != nil {
		log.Println(err)
	}

	the.serverAddress = adr
}

func (the *Network) OpenSocket () {

	socket, err := net.ListenUDP("udp", the.serverAddress)
	if err != nil {
		log.Println(err)
	}

	the.socket = socket
}

func (the *Network) setTimeOutConn (timeOutConnections int) {
	the.timeOutConnections = time.Duration(timeOutConnections) * time.Millisecond
}

func (the *Network) setMessageSize (maxMessageSize int) {
	if maxMessageSize > 1456 {
		// max size udp packet by MTU
		maxMessageSize = 1456
	}
	the.maxMessageSize = maxMessageSize
}



type _ID struct {
	sync.RWMutex
	ID int
}

// Get () int
// safe with Mutex
func (the *_ID) Get() int {
	// block for read
	the.RLock()

	// change id
	the.ID++

	// unblock
	defer the.RUnlock()

	// return result
	return the.ID
}