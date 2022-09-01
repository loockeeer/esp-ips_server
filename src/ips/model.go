package ips

import (
	"github.com/loockeeer/espipsserver/src/tools"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"math"
)

type DistanceRssiModel struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

func NewDistanceRssiModel() DistanceRssiModel {
	return DistanceRssiModel{
		A: 0,
		B: 0,
	}
}

var logpi = 20 * math.Log((4*math.Pi)/0.125) // 0.125 = bluetooth wavelength in m
func (d *DistanceRssiModel) executeLinear(x float64) float64 {
	up := x - d.A - logpi
	down := 10 * d.B
	return up / down
}

func (d *DistanceRssiModel) Execute(x float64) float64 {
	return math.Pow(10, d.executeLinear(x))
}

func (d *DistanceRssiModel) Train(collector RSSICollector, distances map[string]map[string]float64) (err error) {
	errorFunc := func(x []float64) float64 {
		total := 0.0
		m := DistanceRssiModel{
			A: x[0],
			B: x[1],
		}
		for scanner, values := range collector.Data {
			for scanned, rssiqueue := range values {
				if _, ok := distances[scanner]; ok {
					if _, ok := distances[scanner][scanned]; ok {
						avg := float64(tools.Sum[int](tools.Map[TimeEntry, int](rssiqueue.Data, func(v TimeEntry) int {
							return v.RSSI
						}))) / float64(len(rssiqueue.Data))
						total += math.Pow(m.executeLinear(avg)-math.Log10(distances[scanner][scanned]), 2)
					}
				}
			}
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

	result, err := optimize.Minimize(problem, []float64{0, 0}, &optimize.Settings{
		Converger: &optimize.FunctionConverge{
			Iterations: 1000,
		},
	}, &optimize.LBFGS{})

	if err != nil {
		return err
	}
	d.A = result.X[0]
	d.B = result.X[1]
	return nil
}
