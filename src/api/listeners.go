package api

import (
	"espips_server/src/utils"
)

// PositionEvent Push to this signal to emit a device when its position changes
var PositionEvent = &utils.EventEmitter{}

// AppStateChangeEvent Push to this signal to emit when app state is changed
var AppStateChangeEvent = &utils.EventEmitter{}

// DeviceAnnounce Push to this signal to emit when a device is announced
var DeviceAnnounce = &utils.EventEmitter{}
