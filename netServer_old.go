package main

//
//import (
//	"fmt"
//	"log"
//	"net"
//	"time"
//)
//
//type message struct {
//	addr *net.UDPAddr
//	text      []byte
//}
//
//type Messages chan message
//
//type ValidUser map[string]time.Time // *net.UDPAddr
//
//var secret string = "qwerty"
//
//type network struct {
//	secret string
//	maxMessageSize int
//	timeOutConnections time.Duration
//	serverAddress *net.UDPAddr
//	connections map[string]time.Time
//	socket *net.UDPConn
//}
//
//func reader(Conn *net.UDPConn, MessagesFromUser Messages, verifyUser ValidUser) {
//	defer Conn.Close()
//
//	var sms = make([]byte, 1499)
//
//	ticker := time.Tick(8 * time.Millisecond)
//
//	// Conn.SetReadBuffer(10240)
//	for range ticker {
//		size, caddr, err := Conn.ReadFromUDP(sms)
//		if err != nil {
//			log.Println(err)
//		}
//
//		if validation(verifyUser, caddr, &sms, &size); size > 0 {
//			MessagesFromUser <- message{caddr, sms[:size]}
//			// fmt.Println("caddr: ", caddr, "size: ", size, "mess: ", string(sms[:size]))
//		}
//	}
//}
//
//func validation(verifyUser ValidUser, addr *net.UDPAddr, envelope *[]byte, size *int) {
//
//	// check user in verified Users or validation secret
//	if _, ok := verifyUser[addr.String()]; ok || secret == string((*envelope)[:6]) {
//
//		// update time last message
//		verifyUser[addr.String()] = time.Now()
//
//		// del secrets byte from message
//		*envelope = (*envelope)[6:]
//		if string(*envelope) == "ping" {
//			*size = 0
//		}
//	} else {
//		*size = 0
//	}
//}
//
////
//func sender(Conn *net.UDPConn, MessagesToUser Messages) {
//	for envelope := range MessagesToUser {
//		_, err := Conn.WriteTo(envelope.text, envelope.addr)
//		if err != nil {
//			log.Println(err)
//		}
//	}
//}
//
//func handler(verifiedUser ValidUser, MessagesFromUser, MessagesToUser Messages) {
//
//	// slice for offline users
//	offlineUsers := make([]string, 0, 300)
//
//	// timeout for connections // 30sec
//	timeOut := 30000 * time.Millisecond
//
//	for mess := range MessagesFromUser {
//
//		for client, date := range verifiedUser {
//
//			elapsed := time.Now().Sub(date)
//
//			if client != mess.addr.String() && elapsed < timeOut { //
//				addr, err := net.ResolveUDPAddr("udp", client)
//				if err != nil {
//					log.Println(err)
//				}
//				MessagesToUser <- message{addr, mess.text}
//			}
//
//			if elapsed > timeOut {
//				offlineUsers = append(offlineUsers, client)
//			}
//		}
//
//		// remove offline users
//		for _, disconnect := range offlineUsers {
//			delete(verifiedUser, disconnect)
//		}
//		offlineUsers = nil // make([]string, 0, 1000)
//		offlineUsers = make([]string, 0, 300)
//	}
//}
//
//func netStart() {
//	fmt.Println("Server Start")
//
//	// init server
//	sAddr := "0.0.0.0:55442"                     //localhost:0 192.168.0.52  // "0.0.0.0:55442" "192.168.0.52:55442"
//	adr, err := net.ResolveUDPAddr("udp", sAddr) //192.168.0.52:12345
//
//	if err != nil {
//		log.Println(err)
//	}
//
//	Conn, err := net.ListenUDP("udp", adr)
//	if err != nil {
//		log.Println(err)
//	}
//
//	go handler(Users, MessagesFromUsers, MessagesToUsers)
//	go sender(Conn, MessagesToUsers)
//	reader(Conn, MessagesFromUsers, Users)
//}
