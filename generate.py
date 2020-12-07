#!/usr/bin/python
import uuid
import json
import random
import os

states = ['red', 'blue', 'green', 'yellow', 'purple']
records = []

for i in range(20000):
    records.append({'uuid': str(uuid.uuid4()), 'index': i+1, 'state': random.choice(states)})

for record in records:
    print('{"index" : { "_index" : "test", "_id" : "%d" }}' % record['index'])
    print(json.dumps(record))
