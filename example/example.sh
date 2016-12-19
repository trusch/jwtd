#!/bin/bash

docker stop whoami
docker rm whoami
docker run --name whoami -d jwilder/whoami

docker stop jwtd
docker rm jwtd
docker run --name jwtd -d \
  -v $(pwd)/pki/jwtd.crt:/etc/jwtd/jwtd.crt \
  -v $(pwd)/pki/jwtd.key:/etc/jwtd/jwtd.key \
  -v $(pwd)/jwtd.yaml:/etc/jwtd/config.yaml \
  -p 4443:443 \
  trusch/jwtd

docker stop jwtd-proxy
docker rm jwtd-proxy
docker run --name jwtd-proxy -d \
  -v $(pwd)/pki/jwtd.crt:/etc/jwtd-proxy/jwtd.crt \
  -v $(pwd)/jwtd-proxy.yaml:/etc/jwtd-proxy/config.yaml \
  --link whoami \
  -p 8080:8080 \
  trusch/jwtd-proxy

sleep 1

echo "Should not work"
token=$(curl -k --cacert pki/ca.crt --data '{"username":"admin","password":"admin","service":"whoami","labels":{"scope":"user"}}' https://localhost:4443 2>/dev/null)
curl --header "Authorization: Bearer ${token}" http://localhost:8080

echo -en "\nShould work\n"
token=$(curl -k --cacert pki/ca.crt --data '{"username":"admin","password":"admin","service":"whoami","labels":{"scope":"admin"}}' https://localhost:4443 2>/dev/null)
curl --header "Authorization: Bearer ${token}" http://localhost:8080
