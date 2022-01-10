package utils

type EventEmitter struct {
	listeners []chan interface{}
}

func (e *EventEmitter) Listen(listener chan interface{}) {
	e.listeners = append(e.listeners, listener)
}

func (e *EventEmitter) RemoveListener(listener chan interface{}) {
	for i, l := range e.listeners {
		if l == listener {
			e.listeners[i] = e.listeners[len(e.listeners)-1]
			e.listeners = e.listeners[:len(e.listeners)-1]
		}
	}
}

func (e *EventEmitter) Emit(data interface{}) {
	for _, listener := range e.listeners {
		listener <- data
	}
}
