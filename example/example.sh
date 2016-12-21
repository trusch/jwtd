#!/bin/bash

docker stop http-echo
docker rm http-echo
docker run --name http-echo -d trusch/http-echo

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
  -v $(pwd)/pki/http-echo.crt:/etc/jwtd-proxy/http-echo.crt \
  -v $(pwd)/pki/http-echo.key:/etc/jwtd-proxy/http-echo.key \
  -v $(pwd)/jwtd-proxy.yaml:/etc/jwtd-proxy/config.yaml \
  --link http-echo \
  --link jwtd \
  -p 443:443 \
  trusch/jwtd-proxy

sleep 1

echo "get role:admin token as admin"
adminAdminToken=$(curl --cacert pki/ca.crt --data '{"username":"admin","password":"admin","service":"http-echo","labels":{"role":"admin"}}' https://jwtd 2>/dev/null)
echo "get role:user token as admin"
adminUserToken=$(curl --cacert pki/ca.crt --data '{"username":"admin","password":"admin","service":"http-echo","labels":{"role":"user"}}' https://jwtd 2>/dev/null)

echo "get role:admin token as user"
userAdminToken=$(curl --cacert pki/ca.crt --data '{"username":"user","password":"user","service":"http-echo","labels":{"role":"admin"}}' https://jwtd 2>/dev/null)
echo "get role:user token as user"
userUserToken=$(curl --cacert pki/ca.crt --data '{"username":"user","password":"user","service":"http-echo","labels":{"role":"user"}}' https://jwtd 2>/dev/null)

echo "request http-echo /admin with admin-admin token"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${adminAdminToken}" https://http-echo/admin
echo "request http-echo /admin with user-admin token"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${userAdminToken}" https://http-echo/admin

echo "request http-echo /user with admin-user token"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${adminUserToken}" https://http-echo/user
echo "request http-echo /user with user-user token"
curl --cacert pki/ca.crt --header "Authorization: Bearer ${userUserToken}" https://http-echo/user
