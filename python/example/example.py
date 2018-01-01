#!/usr/bin/env python
# -*- coding: utf-8 -*-

from .. import libmqttpy

if __name__ == '__main__':
    client = libmqttpy.Client(server='localhost:1883')
    client.connect()
