package internals

import (
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"math"
)

type Vector []float64

func (p Vector) CustomOptimize(cost func(vector Vector) float64) (newVector Vector, err error) {
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

	result, err := optimize.Minimize(problem, p, nil, &optimize.BFGS{})
	if err != nil {
		return nil, err
	}
	return result.X, nil
}

type Polynomial = Vector

func (p Polynomial) Execute(x float64) float64 {
	result := 0.0
	for i, w := range p {
		result += w * math.Pow(x, float64(i))
	}
	return result
}

func (p Polynomial) Optimize(data map[float64]float64) (poly Polynomial, err error) {
	errorFunc := func(x Vector) float64 {
		result := 0.0
		for key, value := range data {
			result += math.Pow(value-x.Execute(key), 2)
		}
		return result
	}
	vec, err := p.CustomOptimize(errorFunc)
	if err != nil {
		return nil, err
	}
	return vec, nil
}
