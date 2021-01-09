package main

import (
	"fmt"
	"math"
	"math/rand"
)

// the object
type Player struct {
	id           int
	Name         string `json:"name"`
	Coord        Vector `json:"coord"`
	Speed        float64 `json:"speed"`
	status       bool
	targetVector Vector
}

func CreatePlayer(id int, name string, X, Y, speed float64) *Player {
	return &Player{
		id:     id,
		Name:   name,
		Speed:  speed,
		Coord:  Vector{X, Y},
		status: true,
		targetVector: Vector{50,70},
	}
}

func (the *Player) GetName() string {
	return the.Name
}

func (the *Player) GetId() int {
	return the.id
}

func (the *Player) GetCoord() (float64, float64) {
	return the.Coord.X, the.Coord.Y
}

func (the *Player) GetSpeed() float64 {

	return the.Speed
}

// enter world Size {X and Y}
func (the *Player) SetTargetVector(x, y float64) {

	the.targetVector = Vector{Round(x / 2), Round(y / 2)}
}

func (the *Player) SetSpeed(new float64) {

	the.Speed = new
}

func (the *Player) Step() {
	the.move()
	if the.Coord.X >= the.targetVector.X && the.Coord.Y >= the.targetVector.Y {
		rand.Seed(rand.Int63())
		the.targetVector.X = float64(rand.Intn(800))
		the.targetVector.Y = float64(rand.Intn(600))
		fmt.Println("the.targetVector.X: ", the.targetVector.X, " the.targetVector.Y:", the.targetVector.Y)
	}
}

func (the *Player) move() {

	//Speed = math.Sqrt(vec.X*vec.X + vec.Y*vec.Y)

	the.movementCalculation()
}

func (the *Player) getVec() Vector {
	return Vector{
		the.targetVector.X - the.Coord.X,
		the.targetVector.Y - the.Coord.Y,
	}
}

func (the *Player) movementCalculation() {

	var (
		speed float64 = the.Speed
		calcVector Vector
	)

	calcVector.X = the.targetVector.X - the.Coord.X
	calcVector.Y = the.targetVector.Y - the.Coord.Y

	x := math.Pow(the.targetVector.X - the.Coord.X, 2)
	y := math.Pow(the.targetVector.Y - the.Coord.Y, 2)
	vectorLength := math.Sqrt(x+y)

	newVec := Vector{
		calcVector.X / vectorLength,
		calcVector.Y / vectorLength,
	}

	if vectorLength < speed {
		speed = vectorLength
	}

	the.Coord.X = the.Coord.X + (newVec.X * speed)
	the.Coord.Y = the.Coord.Y + (newVec.Y * speed)
	// the.ss = math.Sqrt(vec.X*vec.X + vec.Y*vec.Y) //
	// dot := the.getVec()
	// cos := math.Sqrt(dot.X*dot.X + dot.Y*dot.Y)

	// var CurrVec, NewVec Vector //создаем временные векторы
	// Way, distanceTraveled = normalize(vec)

}
