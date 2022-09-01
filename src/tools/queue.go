package tools

type Queue[T any] struct {
	MaximumSize int
	Data        []T
}

func NewQueue[T any](maxSize int) Queue[T] {
	return Queue[T]{
		MaximumSize: maxSize,
		Data:        []T{},
	}
}

func (q Queue[T]) Add(value T) (popped *T) {
	q.Data = append(q.Data, value)
	if len(q.Data) > q.MaximumSize {
		popped = &q.Data[len(q.Data)-1]
		q.Data = RemoveSlice[T](q.Data, len(q.Data)-1)
		return popped
	}
	return nil
}

func (q Queue[T]) Clear() {
	q.Data = []T{}
}
