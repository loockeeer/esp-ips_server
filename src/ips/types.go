package ips

type Position struct {
	X float64
	Y float64
}

func ListToPosition(list []float64) Position {
	return Position{
		X: list[0],
		Y: list[1],
	}
}
