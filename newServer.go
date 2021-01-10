package main

import (
	"log"
	"net"
	"strings"
	"time"
)

func startConnection(params ServerParameters) (*Network, error) {
	addr, err := net.ResolveUDPAddr(params.networkType, params.address)
	if err != nil {
		return nil, err
	}

	connection, err := net.ListenUDP(params.networkType, addr)
	if err != nil {
		return nil, err
	}

	network := Network{
		secret:             "",
		maxMessageSize:     params.maxMessageSize,
		timeOutConnections: time.Duration(params.timeout) * time.Millisecond,
		serverAddress:      addr,
		socket:             connection,
		status:             true,
		InputMessage:       make(chan InputMessage, 5000),
		systemMessages: systemMessages{
			"ping": []byte("ping"),
			"pong": []byte("pong"),
		},
	}
	return &network, nil
}

func (network *Network) stopConnection() {
	network.socket.Close()
}

func (network *Network) readConnection() {
	defer network.stopConnection()

	buffer := make([]byte, network.maxMessageSize)

	for network.status {
		countBytes, addr, err := network.socket.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
		}

		secret := string(buffer[:len(network.secret)])
		message := buffer[len(network.secret):countBytes]

		if network.secret == strings.ToLower(secret) {
			if _, ok := network.clients[addr.String()]; ok {
				network.clients[addr.String()].lastActivity = time.Now()
			} else {
				newClient := Client{
					lastActivity: time.Now(),
					addr:         addr,
				}
				network.clients[addr.String()] = &newClient
			}

			text := string(message)
			if _, ok := network.systemMessages[text]; !ok {
				network.InputMessage <- InputMessage{addr.String(), message}
			} else {
				if val := checkSystemMessage(text); val != "error" {
					message := []byte(network.secret + val)
					network.Send(addr.String(), message)
				}
			}
		}
	}
}

func (network *Network) Send(addr string, message []byte) {
	network.socket.WriteTo(message, network.clients[addr].addr)
}

func (network *Network) Start() {
	log.Println("UDP server is started: ", network.serverAddress.String())
	network.readConnection()
}

func checkSystemMessage(messageType string) string {
	switch messageType {
	case "ping":
		return "pong"
	case "disconnect":
		return "goodBye"
	default:
		return "error"
	}
}
