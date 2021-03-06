package api

import (
	"espips_server/src/internals"
	"fmt"
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
		"type": &gql.Field{
			Description: "Device type. 0 = station | 1 = beacon",
			Type:        gql.Int,
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
					Type:        gql.String,
					Description: "The address of the device to fetch",
				},
			},
			Resolver: func(ctx gql.Context) (interface{}, error) {
				devices, _ := internals.ListDevices()
				for _, device := range devices {
					if device.Address == ctx.Args()["address"] {
						return internals.GraphQLDevice{
							Address:      *device.Address,
							FriendlyName: *device.FriendlyName,
							X:            device.GetPosition().X,
							Y:            device.GetPosition().Y,
							Speed:        device.GetSpeed(),
							Battery:      device.GetBattery(),
							Type:         int(*device.Type),
						}, nil
					}
				}
				return internals.GraphQLDevice{}, nil
			},
		},
		"devices": &gql.Field{
			Type:        gql.NewList(deviceType),
			Description: "Get device list",
			Resolver: func(ctx gql.Context) (interface{}, error) {
				devices, err := internals.ListDevices()
				if err != nil {
					return nil, err
				}
				if devices == nil {
					return nil, nil
				}
				var data []internals.GraphQLDevice
				for _, device := range devices {
					data = append(data, internals.GraphQLDevice{
						Address:      *device.Address,
						FriendlyName: *device.FriendlyName,
						X:            device.GetPosition().X,
						Y:            device.GetPosition().Y,
						Speed:        device.GetSpeed(),
						Battery:      device.GetBattery(),
						Type:         int(*device.Type),
					})
				}
				return data, nil
			},
		},
	},
}

var subscriptionType = &gql.Object{
	Name: "Subscription",
	Fields: gql.Fields{
		"devices": &gql.Field{
			Type:        deviceType,
			Description: "Subscribe to device data",
			Arguments: gql.Arguments{
				"address": &gql.Argument{
					Type:         gql.String,
					Description:  "The address of the device to subscribe to",
					DefaultValue: nil,
				},
			},
			Resolver: func(ctx gql.Context) (interface{}, error) {
				out := make(chan interface{})
				ch := make(chan interface{})
				PositionEvent.Listen(ch)
				go func() {
					for {
						select {
						case <-ctx.Context().Done():
							PositionEvent.RemoveListener(ch)
							log.Println("Closing a subscription on 'devices'")
							return
						case data := <-ch:
							if address, ok := ctx.Args()["address"]; ok {
								if address == "" || data.(internals.GraphQLDevice).Address == address {
									log.Printf("Sending %#v\n", data)
									out <- data
								}
							} else {
								log.Printf("Sending %#v\n", data)
								out <- data
							}
						}
					}
				}()
				return out, nil
			},
		},
		"app": &gql.Field{
			Type:        gql.Int,
			Description: "Subscribe to app state",
			Resolver: func(context gql.Context) (interface{}, error) {
				ch := make(chan interface{})
				AppStateChangeEvent.Listen(ch)
				go func() {
					for {
						select {
						case <-context.Context().Done():
							AppStateChangeEvent.RemoveListener(ch)
							log.Println("Closing a subscription on 'app'")
							return
						}
					}
				}()
				return ch, nil
			},
		},
		"deviceAnnounce": &gql.Field{
			Type:        deviceType,
			Description: "Subscribe to device announcements",
			Resolver: func(context gql.Context) (interface{}, error) {
				ch := make(chan interface{})
				DeviceAnnounce.Listen(ch)
				go func() {
					for {
						select {
						case <-context.Context().Done():
							DeviceAnnounce.RemoveListener(ch)
							log.Println("Closing a subscription on 'deviceAnnounce")
							return
						}
					}
				}()
				return ch, nil
			},
		},
	},
}

var mutationType = &gql.Object{
	Name: "Mutation",
	Fields: gql.Fields{
		"setMode": &gql.Field{
			Description: "Change app state",
			Arguments: gql.Arguments{
				"mode": &gql.Argument{
					Type:        gql.Int,
					Description: "The mode to set the app to",
				},
			},
			Type: gql.Boolean,
			Resolver: func(context gql.Context) (interface{}, error) {
				internals.AppState = context.Args()["mode"].(int)
				fmt.Println(context.Args()["mode"].(int))
				fmt.Println(internals.AppState)
				GlobalControl(internals.AppState)
				return true, nil
			},
		},
	},
}
