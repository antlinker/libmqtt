#!/usr/bin/python

from distutils.core import setup

setup(
    name = 'libmqtt',
    version = '0.1',
    keywords = ('simple', 'test'),
    description = 'Python mqtt client lib based on libmqtt in go',
    license = 'Apache License v2',

    author = 'goiiot',
    author_email = 'jeffctor@gmail.com',

    py_modules = ['libmqtt'],
    platforms = 'any',
)