# ESP-IPS Server

## Purpose
This is an Indoor Positioning System that relies on the Bluetooth chips of the famous ESP32 microcontrollers.
The idea is to use their capacity to send Bluetooth pings in broadcast mode **while** scanning for other pings, hence building an RSSI (**R**eceived **S**ignal **S**trength **I**ndicator) buffer. The server, hosted in this repository, then collects the data through MQTT and processes the lot.

It features two modes :
- Training : knowing the distance between each ESP32 (beacon and station included), and having a bunch of RSSI between each of them, a linearised model is trained allowing for future distance estimates
- Real Time Positioning : In this mode, beacons are usually placed on moving objects (e.g. : toy cars) and stations fixed at known positions. The server collects all RSSI data between the beacons and the stations, and using the previously trained model estimates every distance and tries to calculate an estimated position for every beacon. This last feature used an optimization algorithm (gradient-descent) to find the best estimate.

Unfortunately, the system never really worked properly. Probably since the antennas used (although straight) weren't made for this purpose, and for the RSSI surely isn't a precise enough metric.

I developped this project in _classe de première_, on the suggestion of my CS lab work teacher in _NSI_ who gave me the opportunity to attend his _BTS_ class' lab work. 

## Instructions

### Config
To start it you will need a .env file like this 
```
# Influx Client/Server Config
INFLUX_HOST= (Optional if no changes are made to docker-compose.yml)
INFLUX_TOKEN=
INFLUX_ORG=
INFLUX_BUCKET=
INFLUX_USERNAME=
INFLUX_PASSWORD=
MODE= (Optional if not in setup)

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
### Setup
For setup, you will need to add a `MODE=setup` entry to the env file
Then you can start influxdb using `docker-compose`:
```shell
docker-compose up influx -d
```
Once it is started, you can stop it
```shell
docker-compose down
```

### Running
Then you can start it using `docker-compose` :
```shell
docker-compose up -d
```

You're done ! The high precision IPS server is now started, and you can enjoy it using the GraphQL API.
