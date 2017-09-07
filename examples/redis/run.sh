#!/usr/bin/env bash

echo "Creating bobcat output..."
mkdir -p output
../../bobcat-darwin -o output/nested_quotes.json ../mongodb/bobcat/nested_quote.lang
../../bobcat-darwin -o output/e.json -s ../mongodb/bobcat/flat_quote.lang
mv output/e-user.json output/users.json
mv output/e-quotes.json output/flat_quotes.json

DOCKER_ID=$(docker ps -f name=redis_server -q)
if [[ -z $DOCKER_ID ]]
then
    echo "Stating redis server container..."
    DOCKER_ID=$(docker run --rm -d --name redis_server kyleolivo/redis-json)
fi
DOCKER_ADDRESS=$(docker inspect $DOCKER_ID | grep -o -m 1 "\"IPAddress\": \".*" | awk -F ": " '{print $2}' | tr -d '",')
DOCKER_PORT="6379"

echo "Loading data into redis..."
docker run --rm -v $(pwd):/data --name python-redis python /bin/bash -c "pip install redis &&
    /data/load_redis.py -a $DOCKER_ADDRESS -p $DOCKER_PORT /data/output/nested_quotes.json &&
    /data/load_redis.py -a $DOCKER_ADDRESS -p $DOCKER_PORT /data/output/users.json &&
    /data/load_redis.py -a $DOCKER_ADDRESS -p $DOCKER_PORT /data/output/flat_quotes.json"