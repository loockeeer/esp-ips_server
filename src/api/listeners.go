package api

import "espips_server/src/internals"

// PositionEmitter Push to this signal to emit a device when its position changes
var PositionEmitter = make(chan internals.GraphQLDevice)

// ChangeAppState Push to this signal to emit when app state is changed
var ChangeAppState = make(chan internals.State)
