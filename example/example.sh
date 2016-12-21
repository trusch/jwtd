#!/bin/bash

docker stop http-echo-1
docker rm http-echo-1
docker run --name http-echo-1 -d trusch/http-echo

docker stop http-echo-2
docker rm http-echo-2
docker run --name http-echo-2 -d trusch/http-echo

docker stop jwtd
docker rm jwtd
docker run --name jwtd -d \
  -v $(pwd)/pki/jwtd.key:/etc/jwtd/jwtd.key \
  -v $(pwd)/jwtd.yaml:/etc/jwtd/config.yaml \
  trusch/jwtd

docker stop jwtd-proxy
docker rm jwtd-proxy
docker run --name jwtd-proxy -d \
  -v $(pwd)/pki/jwtd.crt:/etc/jwtd-proxy/jwtd.crt \
  -v $(pwd)/pki/jwtd.key:/etc/jwtd-proxy/jwtd.key \
  -v $(pwd)/pki/http-echo-1.crt:/etc/jwtd-proxy/http-echo-1.crt \
  -v $(pwd)/pki/http-echo-1.key:/etc/jwtd-proxy/http-echo-1.key \
  -v $(pwd)/pki/http-echo-2.crt:/etc/jwtd-proxy/http-echo-2.crt \
  -v $(pwd)/pki/http-echo-2.key:/etc/jwtd-proxy/http-echo-2.key \
  -v $(pwd)/jwtd-proxy.yaml:/etc/jwtd-proxy/config.yaml \
  --link http-echo-1 \
  --link http-echo-2 \
  --link jwtd \
  -p 443:443 \
  trusch/jwtd-proxy

sleep 1

echo "get admin token"
adminToken=$(curl --cacert pki/ca.crt --data '{"username":"admin","password":"admin","service":"http-echo","labels":{"role":"admin"}}' https://jwtd 2>/dev/null)
echo $adminToken
echo "get user token"
userToken=$(curl --cacert pki/ca.crt --data '{"username":"user","password":"user","service":"http-echo","labels":{"role":"user"}}' https://jwtd 2>/dev/null)
echo $userToken

echo "request http-echo-1 /admin as admin"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${adminToken}" https://http-echo-1/admin
echo "request http-echo-1 /admin as user"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${userToken}" https://http-echo-1/admin
echo "request http-echo-1 /user as admin"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${adminToken}" https://http-echo-1/user
echo "request http-echo-1 /user as user"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${userToken}" https://http-echo-1/user
echo "request http-echo-1 / as nobody"
curl --cacert pki/ca.crt https://http-echo-1/


echo "request http-echo-2 /admin as admin"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${adminToken}" https://http-echo-2/admin
echo "request http-echo-2 /admin as user"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${userToken}" https://http-echo-2/admin
echo "request http-echo-2 /user as admin"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${adminToken}" https://http-echo-2/user
echo "request http-echo-2 /user as user"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${userToken}" https://http-echo-2/user
echo "request http-echo-2 / as nobody"
curl --cacert pki/ca.crt https://http-echo-2/
