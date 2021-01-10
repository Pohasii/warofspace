package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Messages chan InputMessage

type connection map[int]client
// var secret string = "qwerty"

func (the *Network) init () *Network{
	return &Network{
		clients: make(connection, 0),
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
			MessagesFromUser <- InputMessage{caddr, sms[:size]}
		}
	}
}

func (the *Network) validation(userAddr *net.UDPAddr, letter *[]byte, size *int) {

	// check user in verified Users or validation secret
	if _, ok := the.clients[userAddr.String()]; ok || the.secret == string((*letter)[:6]) {

		// update time last InputMessage
		the.clients[userAddr.String()] = time.Now()

		// del secrets byte from InputMessage
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

func (network *Network) OpenSocket () {

	socket, err := net.ListenUDP("udp", network.serverAddress)
	if err != nil {
		log.Println(err)
	}

	network.socket = socket
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
