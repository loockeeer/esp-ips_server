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
	Address      string     `json:"address"`
	FriendlyName string     `json:"friendly_name"`
	X            float64    `json:"x"`
	Y            float64    `json:"y"`
	Type         DeviceType `json:"type"`
}

func (d Device) GetX() float64 {
	result, err := database.Connection.GetPosition(d.Address)
	if err != nil {
		log.Panicln(err)
	}
	return result.X
}

func (d Device) GetY() float64 {
	result, err := database.Connection.GetPosition(d.Address)
	if err != nil {
		log.Panicln(err)
	}
	return result.Y
}

func (d Device) GetSpeed() float64 {
	result, err := database.Connection.GetSpeed(d.Address)
	if err != nil {
		log.Panicln(err)
	}
	return result
}

func (d Device) GetBattery() float64 {
	result, err := database.Connection.GetBattery(d.Address)
	if err != nil {
		log.Panicln(err)
	}
	return result
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
