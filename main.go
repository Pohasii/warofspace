package main

import (
	"fmt"
	"log"
)

var parameters ServerParameters = ServerParameters {
	"",
	"udp",
	30,
	1456,
}

func main() {
	fmt.Println("this is space bro")

	network, err := startConnection(parameters)
	if err != nil {
		log.Fatalln(err)
	}

	go network.Start()

}
