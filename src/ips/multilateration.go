package ips

import (
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"math"
)

func Distance(pos1 Position, pos2 Position) float64 {
	return math.Sqrt(math.Pow(pos1.X-pos2.X, 2) + math.Pow(pos1.Y-pos2.Y, 2))
}

func TrueRangeMultilateration(data map[Position]float64) (pos *Position, err error) {
	errorFunc := func(x []float64) float64 {
		total := 0.0
		for stationPos, dist := range data {
			total += math.Pow(Distance(stationPos, ListToPosition(x))-dist, 2)
		}
		return total
	}
	grad := func(grad, x []float64) {
		fd.Gradient(grad, errorFunc, x, nil)
	}
	hess := func(h *mat.SymDense, x []float64) {
		fd.Hessian(h, errorFunc, x, nil)
	}

	problem := optimize.Problem{
		Func: errorFunc,
		Grad: grad,
		Hess: hess,
	}
	avg := ListToPosition([]float64{0, 0})
	for stationPos, _ := range data {
		avg.X += stationPos.X
		avg.Y += stationPos.Y
	}
	avg.X /= float64(len(data))
	avg.Y /= float64(len(data))

	result, err := optimize.Minimize(problem, []float64{avg.X, avg.Y}, &optimize.Settings{
		Converger: &optimize.FunctionConverge{
			Iterations: 1000,
		},
	}, &optimize.NelderMead{})

	if err != nil {
		return nil, err
	}
	return &Position{
		X: result.X[0],
		Y: result.X[1],
	}, nil
}
