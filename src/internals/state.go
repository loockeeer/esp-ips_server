package internals

type State uint8

const (
	IDLE_STATE State = iota
	RUN_STATE
	INIT_STATE
)

var AppState State = IDLE_STATE
