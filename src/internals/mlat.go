package internals

import (
	"espips_server/src/utils"
	"gonum.org/v1/gonum/optimize"
	"log"
	"math"
)

func GetPosition(distances map[Position]float64) (pos *Position, err error) {
	cost := func(vec Vector) float64 {
		out := 0.0
		for pos, dist := range distances {
			out += math.Pow(utils.Distance(pos.X, pos.Y, vec[0], vec[1]) - dist, 2)
		}
		return out / float64(len(distances))
	}

	vec := Vector{0.0, 0.0}
	for pos, dist := range distances {
		vec[0] += pos.X * 1/dist
		vec[1] += pos.Y * 1/dist
	}
	vec[0] /= float64(len(distances))
	vec[1] /= float64(len(distances))
	vec, err = vec.CustomOptimize(cost, &optimize.NelderMead{}, nil)
	if err != nil {
		log.Panicln(err)
	}
	return &Position{
		X: vec[0],
		Y: vec[1],
	}, err
}
