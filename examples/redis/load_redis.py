#!/usr/bin/env python

import argparse
import json
import redis
from os.path import basename

parser = argparse.ArgumentParser(description='Load entity json files from bobcat into redis')
parser.add_argument('-a', '--address', default="127.0.0.1")
parser.add_argument('-p', '--port', default="6379")
parser.add_argument('FILE')
args = parser.parse_args()

with open(args.FILE) as json_data:
    data = json.load(json_data)

keyname = basename(args.FILE)
r = redis.StrictRedis(host=args.address, port=args.port, db=0)
r.execute_command('JSON.SET', keyname, '.', json.dumps(data))
reply = json.loads(r.execute_command('JSON.GET', keyname))
print("Content stored in redis. Here is the retrieved value:\n\n", reply)