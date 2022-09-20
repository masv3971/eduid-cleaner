#!/usr/bin/env bash

printf "starting Redis...\n"
docker run --rm -it -d -p "6379:6379" -v "${PWD}"/redis-data:/data redis
