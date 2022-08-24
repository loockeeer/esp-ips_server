package ips

import "github.com/loockeeer/espipsserver/src/tools"

type IPSWrapper struct {
	MinimumStationCountForMultilat int
	positionHandler                tools.EventHandler[Position]
	collector                      RSSICollector
}

func NewIPSWrapper(miniStationCount int, maxRssi int) IPSWrapper {
	return IPSWrapper{
		MinimumStationCountForMultilat: miniStationCount,
		positionHandler:                tools.NewEventHandler[Position](),
		collector:                      NewRSSICollector(maxRssi),
	}
}

func (w *IPSWrapper) OnPosition(listener func(*Position)) {
	w.positionHandler.Register(listener)
}

func (w *IPSWrapper) Collect(scanner string, scanned string, rssi int) {
	w.collector.Collect(scanned, scanner, rssi)
	// TODO : Add logic
}
