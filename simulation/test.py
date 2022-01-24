import argparse
import io
import math
import random
import time
import threading

import scipy as scipy
from paho.mqtt import client as mqtt_client
import json
import matplotlib.pyplot as plt

def dist2rssi(distance):
    N = 2
    return -(math.log(distance) * 10 * N)


def dist(x1, y1, x2, y2):
    return math.sqrt((x1 - x2) ** 2 + (y1 - y2) ** 2)


def get_coords(angle):
    return math.cos(angle), math.sin(angle)


class Device:
    def __init__(self, address, friendly_name, x, y, type):
        self.address = address
        self.friendly_name = friendly_name
        self.type = type
        self.x = x
        self.y = y
        self.angle = 0

    @property
    def idle(self):
        return self.type == 3

    @property
    def antenna(self):
        return self.type == 1 or self.type == 2

    @property
    def antennaRun(self):
        return self.type == 2

    @property
    def car(self):
        return self.type == 0

    def loop(self, devices):
        if self.antenna:
            for other in devices:
                if other.address == self.address: continue
                if self.antennaRun and other.car:
                    distance = dist(self.x, self.y, other.x, other.y) + random.uniform(-2, 2)
                    rssi = dist2rssi(distance)
                    yield other.address, rssi, distance

                if self.antenna and not self.antennaRun and other.antenna and not other.antennaRun:
                    distance = dist(self.x, self.y, other.x, other.y) + random.uniform(-2, 2)
                    rssi = dist2rssi(distance)
                    yield other.address, rssi, distance

        elif self.car:
            self.angle += 0.0001
            self.x, self.y = get_coords(self.angle)


def main(devices):
    devices = [Device(device["address"], device["friendly_name"], device["x"], device["y"], 1) if device["type"] == 0 else Device(
        device["address"], device["friendly_name"], 0, 0, 3)
               for device in devices]

    data = []
    for i in range(50):
        for device in devices:
            if not device.antenna: continue
            for oa, rssi, distance in device.loop(devices):
                print(oa, rssi, distance)
                data.append((oa, rssi, distance))

    plt.scatter([item[1] for item in data], [item[2] for item in data])
    plt.show()


if __name__ == "__main__":
    # Load args
    parser = argparse.ArgumentParser(description="Simulates a pseudo real life environment for espips_server")
    parser.add_argument('config', type=str, default="../config/server/config.json")
    args = parser.parse_args()

    # Load devices
    devices = json.loads(io.open(args.config, "r").read())
    main(devices)
