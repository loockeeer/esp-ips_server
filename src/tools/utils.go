package tools

type Number interface {
	int64 | float64
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
