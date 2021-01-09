package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var theID ID

var Users = make(ValidUser, 0)
var MessagesFromUsers = make(Messages, 5000)
var MessagesToUsers = make(Messages, 5000)

func main() {
	fmt.Println("this is space bro")

	var SolarSystem = initTheWorld("Solar", theID.Get(), 800.0, 600.0)
	SolarSystem.addObject(CreateObject(theID.Get(), "Planet", 200, 0.0015, SolarSystem.Size.X/2, SolarSystem.Size.Y/2))
	SolarSystem.addObject(CreateObject(theID.Get(), "PlanetTwo", 110, 0.0030, SolarSystem.Size.X/2, SolarSystem.Size.Y/2))
	SolarSystem.addObject(CreateObject(theID.Get(), "Sun", 0, 0.1, SolarSystem.Size.X/2, SolarSystem.Size.Y/2))
	SolarSystem.addPlayer(CreatePlayer(theID.Get(),"Sonic",0,0, 1.5))

	go netStart()

	var wg = sync.WaitGroup{}

	go func() {
		t := time.Tick(32 * time.Millisecond)
		for range t {
			wg.Add(2)
			go func() {
				defer wg.Done()
				for i, _ := range SolarSystem.ObjectAtMap {
					SolarSystem.ObjectAtMap[i].Step()
				}
			}()

			go func() {
				defer wg.Done()
				for i, _ := range SolarSystem.Players {
					SolarSystem.Players[i].Step()
				}
			}()
			wg.Wait()
		}
	}()

	t := time.Tick(32 * time.Millisecond)

	for range t {

		var addr *net.UDPAddr
		objects, players := SolarSystem.getFrameOfTheWorld()
		// fmt.Println("objects: ",objects, ", players: ",players)
		MessagesFromUsers <- message{addr, []byte("game" + "object" + string(objects) + "\n")}
		MessagesFromUsers <- message{addr, []byte("game" + "player" + string(players) + "\n")}
	}

}
