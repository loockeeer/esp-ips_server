package internals

import (
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
	"math"
)

type Polynomial struct {
	weights []float64
}

func execute(weights []float64, x float64) float64 {
	result := 0.0
	for i, w := range weights {
		result += w * math.Pow(x, float64(i))
	}
	return result
}

func (p Polynomial) Execute(x float64) float64 {
	return execute(p.weights, x)
}

func (p Polynomial) Optimize(data map[float64]float64) (poly Polynomial, err error) {
	errorFunc := func(x []float64) float64 {
		result := 0.0
		for key, value := range data {
			result += math.Pow(key-execute(x, value), 2)
		}
		return result
	}
	grad := func(grad, x []float64) {
		fd.Gradient(grad, errorFunc, x, nil)
	}

	problem := optimize.Problem{
		Func: errorFunc,
		Grad: grad,
	}

	result, err := optimize.Minimize(problem, p.weights, nil, nil)
	if err != nil {
		return Polynomial{}, err
	}
	return Polynomial{
		weights: result.X,
	}, nil
}
