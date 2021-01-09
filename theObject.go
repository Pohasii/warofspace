package main

import "math"

// the object
type SpaceObject struct {
	id           int
	Name         string `json:"name"`
	Coord        Vector `json:"coord"`
	radius       float64
	speed        float64
	currentAngle float64
	status       bool
	worldCenter  Vector
}

func CreateObject(id int, name string, radius, speed float64, worldX, worldY float64) *SpaceObject {
	return &SpaceObject{
		id:     id,
		Name:   name,
		radius: radius,
		speed:  speed,
		worldCenter: Vector{
			worldX,
			worldY,
		},
		status: true,
	}
}

func (the *SpaceObject) getName() string {
	return the.Name
}

func (the *SpaceObject) getId() int {
	return the.id
}

// return X and Y
func (the *SpaceObject) getCoord() (float64, float64) {
	return the.Coord.X, the.Coord.Y
}

// enter world Size {X and Y}
func (the *SpaceObject) SetWorldCenter(x, y float64) {

	// var vx = Math.cos(currentAngle)*radius;
	// var vy = Math.sin(currentAngle)*radius;
	the.worldCenter = Vector{Round(x / 2), Round(y / 2)}
}

func (the *SpaceObject) Step() {
	if the.radius == 0 {
		the.Coord.X = the.worldCenter.X
		the.Coord.Y = the.worldCenter.Y

	} else {
		the.Coord.X = Round(the.worldCenter.X + math.Cos(the.currentAngle)*the.radius)
		the.Coord.Y = Round(the.worldCenter.Y + math.Sin(the.currentAngle)*the.radius)

		// предел для угла
		if the.currentAngle+the.speed >= (math.Pi*2)*the.radius {
			the.currentAngle = 0
		} else {
			the.currentAngle = Round(the.currentAngle + the.speed)
		}
	}
}
