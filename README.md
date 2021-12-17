# ESP-IPS Server

This is the server of a high precision IPS.

## Instructions

To start it you will need a .env file like this 
```
# Influx Client/Server Config
INFLUX_HOST= (Optional if no changes are made to docker-compose.yml)
INFLUX_TOKEN=
INFLUX_ORG=
INFLUX_BUCKET=

# MQTT Client Config
MQTT_PORT=
MQTT_HOST= (Optional if no changes are made to docker-compose.yml)

# API Configuration
API_HOST= 
API_PORT= (Optional if no changes are made to docker-compose.yml)
EXPOSE_PORT=

# General Config
RSSI_BUFFER_SIZE=
INIT_RSSI_BUFFER_SIZE=
RSSI_DISTANCE_ORDER=
```

Then you can start it using `docker-compose` :
```shell
docker-compose up -d
```

You're done ! The high precision IPS server is now started, and you can enjoy it using the GraphQL API.