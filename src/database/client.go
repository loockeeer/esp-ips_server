package database

import (
	"context"
	"espips_server/src/utils"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"time"
)

type Position struct {
	X float64
	Y float64
}

type Database struct {
	client influxdb2.Client
	org    string
	bucket string
}

var Connection *Database
var positionCache map[string]Position
var batteryCache map[string]float64

func Connect(host string, port int, token string, org string, bucket string) *Database {
	log.Println("Connecting with " + token)
	client := influxdb2.NewClient(fmt.Sprintf("http://%s:%d", host, port), token)

	Connection = &Database{
		client,
		org,
		bucket,
	}
	return Connection
}

func (d Database) getWriteClient() api.WriteAPIBlocking {
	return d.client.WriteAPIBlocking(d.org, d.bucket)
}

func (d Database) getReadClient() api.QueryAPI {
	return d.client.QueryAPI(d.org)
}

func (d Database) PushRSSI(scanner string, scannedAddress string, rssi int, trainID string) error {
	point := influxdb2.NewPointWithMeasurement("rssi").
		AddField("scanner", scanner).
		AddField("scanned", scannedAddress).
		AddField("rssi", rssi).
		AddField("trainID", trainID).
		SetTime(time.Now())

	err := d.getWriteClient().WritePoint(context.Background(), point)
	if err != nil {
		return err
	}
	return nil
}

func (d Database) PushPosition(address string, pos Position) error {
	positionCache[address] = pos
	point := influxdb2.NewPointWithMeasurement("position").
		AddTag("address", address).
		AddField("x", pos.X).
		AddField("y", pos.Y).
		SetTime(time.Now())
	err := d.getWriteClient().WritePoint(context.Background(), point)
	if err != nil {
		return err
	}
	return nil
}

func (d Database) PushBattery(address string, battery float64) error {
	batteryCache[address] = battery
	point := influxdb2.NewPointWithMeasurement("battery").
		AddTag("address", address).
		AddField("battery", battery).
		SetTime(time.Now())
	err := d.getWriteClient().WritePoint(context.Background(), point)
	if err != nil {
		return err
	}
	return nil
}

func (d Database) GetSpeed(address string) (speed float64, err error) {
	result, err := d.getReadClient().Query(context.Background(), fmt.Sprintf(`
from(bucket:"%s")
  |> range(start: -3s)
  |> filter(fn: (r) =>
    r._measurement == "position" and
	r.address == "%s"
  )
  |> yield()
`, d.bucket, address))
	if err != nil {
		return -1.0, err
	}
	i := 0
	startX := 0.0
	startY := 0.0
	x := 0.0
	y := 0.0
	for result.Next() {
		x = result.Record().ValueByKey("x").(float64)
		y = result.Record().ValueByKey("y").(float64)
		if i == 0 {
			startX = x
			startY = y
		}
		i++
	}
	return utils.Distance(startX, startY, x, y) / 3, nil
}

func (d Database) GetPosition(address string) (position *Position, err error) {
	if val, ok := positionCache[address]; ok {
		return &val, nil
	}
	result, err := d.getReadClient().Query(context.Background(), fmt.Sprintf(`
from(bucket:"%s")
  |> range(
	start: -10s,
	stop: now()
	)
  |> limit(n: 1)
  |> filter(fn: (r) =>
    r._measurement == "position" and
	r.address == "%s"
  )
  |> yield()
`, d.bucket, address))
	if err != nil {
		return nil, err
	} else {
		if result.Next() {
			return &Position{
				X: result.Record().ValueByKey("x").(float64),
				Y: result.Record().ValueByKey("y").(float64),
			}, nil
		} else {
			return nil, nil
		}
	}
}

func (d Database) GetBattery(address string) (batteryLevel float64, err error) {
	if val, ok := batteryCache[address]; ok {
		return val, nil
	}
	result, err := d.getReadClient().Query(context.Background(), fmt.Sprintf(`
from(bucket:"%s")
  |> range(
	start: -10s,
	stop: now()
	)
  |> limit(n: 1)
  |> filter(fn: (r) =>
    r._measurement == "battery" and
	r.address == "%s"
  )
  |> yield()
`, d.bucket, address))
	if err != nil {
		return -1.0, err
	} else {
		if result.Next() {
			return result.Record().ValueByKey("battery").(float64), nil
		} else {
			return -1.0, nil
		}
	}
}

func (d Database) GetRSSIHistory(trainID string) (history map[string]map[string][]float64, err error) {
	result, err := d.getReadClient().Query(context.Background(), fmt.Sprintf(`
from(bucket:"%s")
  |> range(
	start: -15m,
	stop: now()
	) 
  |> filter(fn: (r) =>
    r._measurement == "rssi" and
	r.trainID = "%s"
  )
  |> group(columns: ["scanner", "scanned"]
  |> yield()
`, d.bucket, trainID))
	if err != nil {
		return nil, err
	} else {
		for result.Next() {
			scanner := result.TableMetadata().Column(0).DefaultValue()
			scanned := result.TableMetadata().Column(1).DefaultValue()
			if result.TableChanged() {
				history[scanner][scanned] = []float64{}
			}
			history[scanner][scanned] = append(history[scanner][scanned], result.Record().ValueByKey("rssi").(float64))
		}
		return history, nil
	}
}
