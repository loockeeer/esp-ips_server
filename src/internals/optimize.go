package internals

import (
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
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
		return Vector{}, err
	}
	return result.X, nil
}
