package ips

import (
	"github.com/loockeeer/espipsserver/src/tools"
	"time"
)

type IPSWrapper struct {
	MinimumStationCountForMultilat int
	Model                          DistanceRssiModel
	Positions                      map[string]Position
	positionHandler                tools.EventHandler[PositionEvent]
	collector                      RSSICollector
}

type PositionEvent struct {
	Device   string
	Position Position
}

func NewIPSWrapper(minStationCount int, maxRssi int, maxTimeDiff time.Duration, model DistanceRssiModel, pos map[string]Position) IPSWrapper {
	return IPSWrapper{
		MinimumStationCountForMultilat: minStationCount,
		Model:                          model,
		Positions:                      pos,
		positionHandler:                tools.NewEventHandler[PositionEvent](),
		collector:                      NewRSSICollector(maxRssi, maxTimeDiff),
	}
}

func (w *IPSWrapper) OnPosition(listener func(*PositionEvent)) {
	w.positionHandler.Register(listener)
}

func (w *IPSWrapper) Collect(scanner string, scanned string, rssi int) error {
	w.collector.Collect(scanned, scanner, rssi)
	var data map[Position]float64
	for scannerAddress, queue := range w.collector.Data[scanned] {
		if len(queue.Data) == queue.MaximumSize {
			data[w.Positions[scannerAddress]] = w.Model.Execute(float64(tools.Sum[int](tools.Map[TimeEntry, int](queue.Data, func(v TimeEntry) int {
				return v.RSSI
			}))) / float64(len(queue.Data)))
		}
	}
	if len(data) >= w.MinimumStationCountForMultilat {
		pos, err := TrueRangeMultilateration(data)
		w.positionHandler.Dispatch(&PositionEvent{
			Position: *pos,
			Device:   scanned,
		})
		return err
	}
	return nil
}
