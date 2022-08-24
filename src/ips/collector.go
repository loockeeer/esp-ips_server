package ips

import "github.com/loockeeer/espipsserver/src/tools"

type RSSICollector struct {
	MaximumRSSIValues int
	Data              map[string]map[string]tools.Queue[int]
}

func NewRSSICollector(maxSize int) RSSICollector {
	return RSSICollector{
		MaximumRSSIValues: maxSize,
		Data:              map[string]map[string]tools.Queue[int]{},
	}
}

// Collect TODO: Refactor to make it prettier (URGENT)
func (c *RSSICollector) Collect(scanner string, scanned string, rssi int) {
	if _, ok := c.Data[scanner]; ok {
		if _, ok2 := c.Data[scanner][scanned]; ok2 {
			c.Data[scanner][scanned].Add(rssi)
		} else {
			c.Data[scanner][scanned] = tools.NewQueue[int](c.MaximumRSSIValues)
			c.Data[scanner][scanned].Add(rssi)
		}
	} else {
		c.Data[scanner] = map[string]tools.Queue[int]{}
		c.Data[scanner][scanned] = tools.NewQueue[int](c.MaximumRSSIValues)
		c.Data[scanner][scanned].Add(rssi)
	}
}
