package internals

import (
	"espips_server/src/database"
	"log"
)

type DeviceType int

const (
	AntennaType DeviceType = iota
	CarType
)

type Device struct {
	Address      *string     `json:"address"`
	FriendlyName *string     `json:"friendly_name"`
	X            *float64    `json:"x"`
	Y            *float64    `json:"y"`
	Type         *DeviceType `json:"type"`
}

func (d Device) GetPosition() (position *Position) {
	if d.X != nil && d.Y != nil {
		return &Position{X: *d.X, Y: *d.Y}
	}
	pos, _ := database.Connection.GetPosition(*d.Address)
	return &Position{X: pos.X, Y: pos.Y}
}

func (d Device) GetSpeed() float64 {
	result, err := database.Connection.GetSpeed(*d.Address)
	if err != nil {
		log.Panicln(err)
	}
	return result
}

func (d Device) GetBattery() float64 {
	result, err := database.Connection.GetBattery(*d.Address)
	if err != nil {
		log.Panicln(err)
	}
	return result
}

func (d Device) SetPosition(pos *Position) error {
	return database.Connection.PushPosition(*d.Address, database.Position{X: pos.X, Y: pos.Y})
}

func (d Device) SetBattery(battery float64) error {
	return database.Connection.PushBattery(*d.Address, battery)
}

type GraphQLDevice struct {
	Address      string  `json:"address"`
	FriendlyName string  `json:"friendlyName"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Speed        float64 `json:"speed"`
	Battery      float64 `json:"battery"`
	Type         int     `json:"type"`
}

type Position struct {
	X float64
	Y float64
}
