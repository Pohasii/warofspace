package main

import (
	"encoding/json"
	"log"
)

// the world map
type WorldObject struct {
	Name string
	ID   int
	Size struct {
		X float64
		Y float64
	}
	ObjectAtMap map[int]*SpaceObject `json:"set1"`
	Players map[int]*Player `json:"set2"`
	Status      bool
}

func initTheWorld(name string, id int, x float64, y float64) *WorldObject {
	return &WorldObject{
		Name: name,
		ID:   id,
		Size: struct {
			X float64
			Y float64
		}{X: x, Y: y},
		ObjectAtMap: make(map[int]*SpaceObject),
		Players: make(map[int]*Player),
		Status:      true,
	}
}

func (the *WorldObject) getName() string {
	return the.Name
}

func (the *WorldObject) getID() int {
	return the.ID
}

func (the *WorldObject) getStatus() bool {
	return the.Status
}

func (the *WorldObject) getSize() (float64, float64) {
	return the.Size.X, the.Size.Y
}

func (the *WorldObject) checkExistingObject(id int) bool {

	object, ok := the.ObjectAtMap[id]

	// check type of value
	var t interface{} = object
	_, check := t.(SpaceObject)

	if check {
		return ok
	}

	return ok
}

func (the *WorldObject) addObject(new *SpaceObject) bool {
	if !the.checkExistingObject(new.getId()) {
		the.ObjectAtMap[new.getId()] = new
	}

	return the.checkExistingObject(new.getId())
}

func (the *WorldObject) deleteObjById(id int) bool {
	delete(the.ObjectAtMap, id)
	return !the.checkExistingObject(id)
}

func (the *WorldObject) getFrameOfTheWorld() (objects []byte, players []byte) {

	objects, err := json.Marshal(&the.ObjectAtMap)
	if err != nil {
		log.Println(err)
	}

	p := (*the).Players
	// fmt.Println("p ", p, "p ", &p)
	players, err = json.Marshal(&p)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(string(val))
	return objects, players
}

func (the *WorldObject) addPlayer(new *Player) bool {
	if !the.checkExistingPlayer(new.GetId()) {
		the.Players[new.GetId()] = new
	}

	return the.checkExistingPlayer(new.GetId())
}

func (the *WorldObject) checkExistingPlayer (id int) bool {

	player, ok := the.Players[id]

	// check type of value
	var t interface{} = player
	_, check := t.(Player)

	if check {
		return ok
	}

	return ok
}