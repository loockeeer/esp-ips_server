package internals

const (
	IDLE_STATE = iota
	RUN_STATE
	INIT_STATE
)

var AppState = IDLE_STATE
