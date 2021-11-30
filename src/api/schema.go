package api

import (
	"espips_server/src/internals"
	"github.com/rigglo/gql"
	"log"
)

var deviceType = &gql.Object{
	Name: "Device",
	Fields: gql.Fields{
		"address": &gql.Field{
			Description: "The address of the device",
			Type:        gql.String,
		},
		"friendlyName": &gql.Field{
			Description: "Friendly name of device defined by server",
			Type:        gql.String,
		},
		"x": &gql.Field{
			Description: "x position of device",
			Type:        gql.Float,
		},
		"y": &gql.Field{
			Description: "y position of device",
			Type:        gql.Float,
		},
		"speed": &gql.Field{
			Description: "Speed of device in m/s",
			Type:        gql.Float,
		},
		"battery": &gql.Field{
			Description: "Device battery in V",
			Type:        gql.Float,
		},
	},
}

var queryType = &gql.Object{
	Name: "Query",
	Fields: gql.Fields{
		"device": &gql.Field{
			Type:        deviceType,
			Description: "Get device by address",

			Arguments: gql.Arguments{
				"address": &gql.Argument{
					Type: gql.String,
				},
			},
			Resolver: func(ctx gql.Context) (interface{}, error) {
				return nil, nil
			},
		},
		"devices": &gql.Field{
			Type:        gql.NewList(deviceType),
			Description: "Get device list",
			Resolver: func(ctx gql.Context) (interface{}, error) {
				return nil, nil
			},
		},
	},
}

var subscriptionsType = &gql.Object{
	Name: "Subscriptions",
	Fields: gql.Fields{
		"device": &gql.Field{
			Type:        deviceType,
			Description: "Subscribe to device data",
			Arguments: gql.Arguments{
				"address": &gql.Argument{
					Type: gql.String,
				},
			},
			Resolver: func(ctx gql.Context) (interface{}, error) {
				out := make(chan internals.Device)
				go func() {
					for {
						select {
						case <-ctx.Context().Done():
							log.Println("Closing a connection")
							return
						case data := <-PositionEmitter:
							if data.Address == ctx.Args()["address"] {
								log.Printf("Sending %#v\n", data)
								out <- data
							}
						}
					}
				}()
				return out, nil
			},
		},
	},
}
