package tools

type Number interface {
	int | int8 | int16 | int32 | int64 | float64 | float32
}

func RemoveSlice[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func Map[T any, V any](s []T, f func(T) V) []V {
	var ns []V
	for _, v := range s {
		ns = append(ns, f(v))
	}
	return ns
}

func Sum[T Number](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

func Find[T any](s []T, f func(T)bool) *T {
	for _, v := range s {
		if f(v) {
			return &v
		}
	}
	return nil
}