#!/bin/bash

docker push trusch/jwtd:latest
docker push trusch/jwtd-proxy:latest

docker tag trusch/jwtd:latest trusch/jwtd:$(git describe)
docker tag trusch/jwtd-proxy:latest trusch/jwtd-proxy:$(git describe)

docker push trusch/jwtd:$(git describe)
docker push trusch/jwtd-proxy:$(git describe)

exit 0
