package ips

import (
	"github.com/loockeeer/espipsserver/src/tools"
	"time"
)

type TimeEntry struct {
	RSSI int
	Time time.Time
}

type RSSICollector struct {
	MaximumRSSIValues     int
	MaximumTimeDifference time.Duration
	Data                  map[string]map[string]tools.Queue[TimeEntry]
}

func NewRSSICollector(maxSize int, maxTimeDiff time.Duration) RSSICollector {
	return RSSICollector{
		MaximumRSSIValues:     maxSize,
		MaximumTimeDifference: maxTimeDiff,
		Data:                  map[string]map[string]tools.Queue[TimeEntry]{},
	}
}

// Collect TODO: Refactor to make it prettier (URGENT)
func (c *RSSICollector) Collect(scanner string, scanned string, rssi int) {
	if _, ok := c.Data[scanner]; ok {
		if _, ok2 := c.Data[scanner][scanned]; ok2 {
			c.Data[scanner][scanned].Add(TimeEntry{rssi, time.Now()})
		} else {
			c.Data[scanner][scanned] = tools.NewQueue[TimeEntry](c.MaximumRSSIValues)
			c.Data[scanner][scanned].Add(TimeEntry{rssi, time.Now()})
		}
	} else {
		c.Data[scanner] = map[string]tools.Queue[TimeEntry]{}
		c.Data[scanner][scanned] = tools.NewQueue[TimeEntry](c.MaximumRSSIValues)
		c.Data[scanner][scanned].Add(TimeEntry{rssi, time.Now()})
	}
	var nearest time.Time
	for d1, _ := range c.Data {
		for _, queue := range c.Data[d1] {
			for _, r := range queue.Data {
				if r.Time.After(nearest) {
					nearest = r.Time
				}
			}
		}
	}
	for d1, _ := range c.Data {
		for _, queue := range c.Data[d1] {
			for i, r := range queue.Data {
				if r.Time.Sub(nearest) > c.MaximumTimeDifference {
					queue.Data = tools.RemoveSlice[TimeEntry](queue.Data, i)
				}
			}
		}
	}
}
