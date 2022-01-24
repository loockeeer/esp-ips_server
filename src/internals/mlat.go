package internals

import (
	"espips_server/src/utils"
	"gonum.org/v1/gonum/optimize"
	"log"
)

func GetPosition(distances map[Position]float64) (pos *Position, err error) {
	cost := func(vec Vector) float64 {
		out := 0.0
		for pos, dist := range distances {
			out += utils.Distance(pos.X, pos.Y, vec[0], vec[1]) - dist
		}
		return out / float64(len(distances))
	}

	vec := Vector{0.0, 0.0}
	vec, err = vec.CustomOptimize(cost, &optimize.GradientDescent{}, nil)
	if err != nil {
		log.Panicln(err)
	}
	return &Position{
		X: vec[0],
		Y: vec[1],
	}, err
}
