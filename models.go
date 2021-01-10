package main

import (
	"net"
	"sync"
	"time"
)

type ServerParameters struct {
	address string
	networkType string
	timeout int
	maxMessageSize int
}

type Client struct {
	ID int
	lastActivity time.Time
	addr string
}

type Network struct {
	nextClientID       _ID
	secret             string
	maxMessageSize     int
	timeOutConnections time.Duration
	serverAddress      *net.UDPAddr
	clients            map[string]*Client
	socket             *net.UDPConn
	status             bool
	inputMessage       chan InputMessage
	outputMessage      chan InputMessage
}

type InputMessage struct {
	_ID int
	text []byte
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