# ESP-IPS Server

This is the server of a high precision IPS.

## Instructions

To start it you will need a .env file like this 
```
INFLUX_TOKEN=
INFLUX_ORG=
INFLUX_BUCKET=
MQTT_PORT=
API_HOST=
API_PORT=
```

Then you can start it using `docker-compose` :
```shell
docker-compose up -d
```

You're done ! The high precision IPS server is now started, and you can enjoy it using the GraphQL API.