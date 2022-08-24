package tools

import "fmt"

type EventHandler[T any] struct {
	listeners []func(event *T)
}

func NewEventHandler[T any]() EventHandler[T] {
	return EventHandler[T]{
		listeners: []func(event *T){},
	}
}

func (h *EventHandler[T]) Register(listener func(*T)) {
	h.listeners = append(h.listeners, listener)
}

func (h *EventHandler[T]) Dispatch(event *T) {
	for _, listener := range h.listeners {
		go listener(event)
	}
}

func (h *EventHandler[T]) Remove(listener func(T)) {
	for i, listen := range h.listeners {
		if fmt.Sprintf("%V", listener) == fmt.Sprintf("%V", listen) {
			h.listeners = RemoveSlice[func(event *T)](h.listeners, i)
		}
	}
}
