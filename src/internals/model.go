package internals

import (
	"fmt"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"math"
	"math/rand"
)

type Vector []float64

func (p Vector) CustomOptimize(cost func(vector Vector) float64, method optimize.Method, settings *optimize.Settings) (newVector Vector, err error) {
	errorFunc := func(x []float64) float64 {
		return cost(x)
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
	if settings == nil {
		settings = &optimize.Settings{
			Converger: &optimize.FunctionConverge{
				Iterations: 1000,
			},
		}
	}

	result, err := optimize.Minimize(problem, p, settings, method)
	if err != nil {
		return nil, err
	}
	return result.X, nil
}

var DistanceRssi = CreateModel()

type Model = Vector

func CreateModel() Model {
	return Model{10 * rand.Float64(), 10 * rand.Float64()}
}

var logpi = 20 * math.Log((4*math.Pi)/0.125) // 0.125 = bluetooth wavelength in m
func (m Model) Execute(x float64) float64 {
	up := x - logpi
	down := 10 * m[0]
	return math.Pow(10, up/down) + m[1]
}

func (m Model) Optimize(data map[float64]float64) (model Model, err error) {
	errorFunc := func(x Vector) float64 {
		result := 0.0
		for key, value := range data {
			result += math.Abs(value - x.Execute(key))
		}
		fmt.Println(result / float64(len(data)))
		return result / float64(len(data))
	}
	vec, err := m.CustomOptimize(errorFunc, &optimize.NelderMead{}, nil)
	if err != nil {
		return nil, err
	}
	return vec, nil
}
