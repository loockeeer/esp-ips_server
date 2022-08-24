package tools

type Queue[T any] struct {
	MaximumSize int
	data        []T
}

func NewQueue[T any](maxSize int) Queue[T] {
	return Queue[T]{
		MaximumSize: maxSize,
		data:        []T{},
	}
}

func (q Queue[T]) Add(value T) (popped *T) {
	q.data = append(q.data, value)
	if len(q.data) > q.MaximumSize {
		popped = &q.data[len(q.data)-1]
		q.data = RemoveSlice[T](q.data, len(q.data)-1)
		return popped
	}
	return nil
}

func (q Queue[T]) Get() []T {
	return q.data
}

func (q Queue[T]) Clear() {
	q.data = []T{}
}
