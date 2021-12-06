package internals

import (
	"math"
)

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
		return Polynomial{}, err
	}
	return vec, nil
}
