package main

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func Round(number float64) float64 {
	return number //math.Floor(number*100) / 100
}
