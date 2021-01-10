package main

import (
	"log"
	"net"
	"strings"
	"time"
)

func startConnection (params ServerParameters) (*Network, error)  {
	addr, err := net.ResolveUDPAddr(params.networkType, params.address)
	if err != nil {
		return nil, err
	}

	connection, err := net.ListenUDP(params.networkType, addr)
	if err != nil {
		return nil, err
	}

	network := Network{
		secret: "",
		maxMessageSize: params.maxMessageSize,
		timeOutConnections: time.Duration(params.timeout) * time.Millisecond,
		serverAddress: addr,
		socket: connection,
		status: true,
		inputMessage: make(chan InputMessage, 5000),
		outputMessage: make(chan InputMessage, 5000),
	}
	return &network, nil
}

func (network *Network) stopConnection () {
	network.socket.Close()
}

func (network *Network) readConnection () {
	defer network.stopConnection()

	buffer := make([]byte, network.maxMessageSize)

	for network.status {
		countBytes, addr, err := network.socket.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
		}

		secret := string(buffer[:len(network.secret)])
		message := buffer[len(network.secret): countBytes]

		if network.secret == strings.ToLower(secret) {
			if client, ok := network.clients[addr.String()]; ok {
				network.clients[addr.String()].lastActivity = time.Now()
				network.inputMessage <- InputMessage{ client.ID, message}
			} else {
				newClient := Client{
					ID: network.nextClientID.Get(),
					lastActivity: time.Now(),
					addr: addr.String(),
				}
				network.clients[addr.String()] = &newClient
				network.inputMessage <- InputMessage{ newClient.ID, message}
			}
		}
	}

}