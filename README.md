# ESP-IPS Server

This is an implementation of an IPS server usable with the esp_ips client (available on my github)

## Building

```bash
make build
```

You will find the binary in build/

## Configuration

### Device list and model

This should be a JSON file in this format :

```json
{
  "devices": [
    {
      "address": "", // The mac address of the device
      "friendlyName": "", // A name easy to remember (not used by the server)
      "type": "", // Either station or beacon
      "x": 0.0, // Only if station
      "y": 0.0 // Only if station 
    }
  ],
  "model": { // This will be created by the server once started in train mode
    "a": 3.4313,  // Sample value
    "b": 1.1241   // Sample value
  }
}
```

### Environment

```dotenv
CONFIG_FILE=/some/config/file/path.json
MODE=train|run
RSSI_MIN_TRAIN=100 // Minimum amount of RSSI values per device pair for training
RSSI_MAX_RUN=3 // Minimum and maximum amount of RSSI values needed per device pair for running
MAX_TIME_DIFFERENCE=10s // TTL for cache values, relative to the last RSSI value recorded
MQTT_HOST=127.0.0.1
MQTT_PORT=1883
```

## Instructions

You will need two things to make it work :
- An MQTT server
- Some ESP-32 flashed with my esp-ips

Then simply launch the binary with the correct env


## Note

The server is designed to be used with docker, and more specifically docker-compose, hence the use of environment variables
for handling the config part instead of regular CLI arguments (It's way more convenient)