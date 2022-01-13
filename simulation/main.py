import argparse
import io
import math
import random
import time
import threading
from paho.mqtt import client as mqtt_client
import json

def dist2rssi(distance):
    N = 2
    return -(math.log(distance)*10*N)

def dist(x1, y1, x2, y2):
    return math.sqrt((x1 - x2) ** 2 + (y1 - y2) ** 2)

def get_coords(angle):
    return math.cos(angle), math.sin(angle)

class Device:
    def __init__(self, address, friendly_name, x, y):
        self.address = address
        self.friendly_name = friendly_name
        self.type = 3
        self.x = x
        self.y = y
        self.angle = 0

    @property
    def idle(self):
        return self.type == 3

    @property
    def station(self):
        return self.type == 1 or self.type == 2

    @property
    def stationRun(self):
        return self.type == 2

    @property
    def beacon(self):
        return self.type == 0

    def loop(self, devices, client):
        if self.station:
            for other in devices:
                if other.address == self.address: continue
                if self.stationRun and other.beacon:
                    distance = dist(self.x, self.y, other.x, other.y) + random.randrange(-2, 2)
                    rssi = dist2rssi(distance)
                    print(self.address, other.address, rssi, distance)
                    client.publish(f"rssi/{self.address}", f"{other.address},{str(int(rssi))}")

                if self.station and not self.stationRun and other.station and not other.stationRun:
                    distance = dist(self.x, self.y, other.x, other.y) + random.randrange(-2, 2)
                    rssi = dist2rssi(distance)
                    print(self.address, other.address, rssi, distance)
                    client.publish(f"rssi/{self.address}", f"{other.address},{str(int(rssi))}")
        elif self.beacon:
            self.angle += 0.001
            self.x, self.y = get_coords(self.angle)

    def ack(self, client):
        if not self.station: return
        client.publish(f"cc/{self.address}", "4")

def main(mqtt_host, mqtt_port, devices):
    def on_message(client, userdata, message):
        content, topic = message.payload.decode(), message.topic
        print(content, topic)
        if topic.startswith("cc"):
            address = topic.split("/")[1]
            for device in devices:
                if device.address == address:
                    if content == "4": continue
                    print(f"Device {device.friendly_name} switched to {content}")
                    device.type = int(content)
                    #device.ack(client)

    def on_connect(client, userdata, flags, rc):
        print("MQTT Connected")
        client.subscribe("cc/+")
        client.subscribe("cc")

    client = mqtt_client.Client("simulation")
    client.on_message = on_message
    client.on_connect = on_connect

    client.connect(mqtt_host, mqtt_port)

    devices = [Device(device["address"], device["friendly_name"], device["x"], device["y"]) if device["type"] == 0 else Device(device["address"], device["friendly_name"], 0, 0)
               for device in devices]

    for device in devices:
        device.ack(client)

    threading.Thread(target=client.loop_forever).start()
    print("loop in thread")
    while True:
        for device in devices:
            time.sleep(0.1)
            device.loop(devices, client)


if __name__ == "__main__":
    # Load args
    parser = argparse.ArgumentParser(description="Simulates a pseudo real life environment for espips_server")
    parser.add_argument('mqtt_host', type=str)
    parser.add_argument('mqtt_port', type=int)
    parser.add_argument('config', type=str, default="../config/server/config.json")
    args = parser.parse_args()

    # Load devices
    devices = json.loads(io.open(args.config, "r").read())
    main(args.mqtt_host, args.mqtt_port, devices)
