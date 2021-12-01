package database

import (
	"context"
	"espips_server/src/internals"
	"espips_server/src/utils"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"time"
)

type Database struct {
	client influxdb2.Client
	org    string
	bucket string
}

var Connection Database

func Connect(host string, port int, token string, org string, bucket string) Database {
	client := influxdb2.NewClient(fmt.Sprintf("http://%s:%d²", host, port), token)

	Connection = Database{
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

func (d Database) PushKeep(address string) error {
	point := influxdb2.NewPointWithMeasurement("keep").
		AddField("address", address).
		SetTime(time.Now())

	err := d.getWriteClient().WritePoint(context.Background(), point)
	if err != nil {
		return err
	}
	return nil
}

func (d Database) PushPosition(address string, x float64, y float64) error {
	point := influxdb2.NewPointWithMeasurement("position").
		AddTag("address", address).
		AddField("x", x).
		AddField("y", y).
		SetTime(time.Now())
	err := d.getWriteClient().WritePoint(context.Background(), point)
	if err != nil {
		return err
	}
	return nil
}

func (d Database) PushBattery(address string, battery float64) error {
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

func (d Database) GetPosition(address string) (position internals.Position, err error) {
	result, err := d.getReadClient().Query(context.Background(), fmt.Sprintf(`
from(bucket:"%s")
  |> limit(n: 1)
  |> filter(fn: (r) =>
    r._measurement == "position" and
	r.address == "%s"
  )
  |> yield()
`, d.bucket, address))
	if err != nil {
		return internals.Position{}, err
	} else {
		result.Next()
		return internals.Position{
			X: result.Record().ValueByKey("x").(float64),
			Y: result.Record().ValueByKey("y").(float64),
		}, nil
	}
}

func (d Database) GetBattery(address string) (batteryLevel float64, err error) {
	result, err := d.getReadClient().Query(context.Background(), fmt.Sprintf(`
from(bucket:"%s")
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
		result.Next()
		return result.Record().ValueByKey("battery").(float64), nil
	}
}
