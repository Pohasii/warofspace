package main

import (
	"net"
	"time"
)

type ServerParameters struct {
	address string
	networkType string
	timeout int
	maxMessageSize int
}

type Client struct {
	lastActivity time.Time
	addr *net.UDPAddr
}

type Network struct {
	secret             string
	maxMessageSize     int
	timeOutConnections time.Duration
	serverAddress      *net.UDPAddr
	clients            map[string]*Client
	socket             *net.UDPConn
	status             bool
	InputMessage       chan InputMessage
	systemMessages systemMessages
}

type InputMessage struct {
	Addr string
	Text []byte
}

type systemMessages map[string][]byte
//type _ID struct {
//	sync.RWMutex
//	ID int
//}
//
//// Get () int
//// safe with Mutex
//func (the *_ID) Get() int {
//	// block for read
//	the.RLock()
//
//	// change id
//	the.ID++
//
//	// unblock
//	defer the.RUnlock()
//
//	// return result
//	return the.ID
//}