#!/bin/bash

function setupDocker {
  cp jwtd.yaml.tmpl jwtd.yaml

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
}

function getToken {
  USER=$1
  PASSWORD=$2
  SERVICE=$3
  LABELKEY=$4
  LABELVALUE=$5
  payload=$(printf '{"username":"%s","password":"%s","service":"%s","labels":{"%s":"%s"}}' $USER $PASSWORD $SERVICE $LABELKEY $LABELVALUE)
  curl --cacert pki/ca.crt --data $payload https://jwtd/token 2>/dev/null
}

function doGetRequest {
  TOKEN=$1
  URL=$2
  echo $(curl --cacert pki/ca.crt --header "Authorization: Bearer ${TOKEN}" $URL 2>/dev/null)
}

function doPostRequest {
  TOKEN=$1
  URL=$2
  DATA=$3
  echo $(curl --cacert pki/ca.crt -X POST --data "${DATA}" --header "Authorization: Bearer ${TOKEN}" $URL 2>/dev/null)
}

function doPatchRequest {
  TOKEN=$1
  URL=$2
  DATA=$3
  echo $(curl --cacert pki/ca.crt -X PATCH --data "${DATA}" --header "Authorization: Bearer ${TOKEN}" $URL 2>/dev/null)
}

function setupUserAndGroup {
  jwtdAdminToken=$(getToken admin admin jwtd role admin)
  echo "get user list"
  doGetRequest $jwtdAdminToken https://jwtd/project/default/user
  echo "create user 'example'"
  doPostRequest $jwtdAdminToken https://jwtd/project/default/user '{"username":"example","password":"example"}'
  echo "get user 'example'"
  doGetRequest $jwtdAdminToken https://jwtd/project/default/user/example
  echo "create group 'http-echo-admin'"
  doPostRequest $jwtdAdminToken https://jwtd/project/default/group '{"name":"http-echo-admin","rights":{"http-echo":{"role":"admin"}}}'
  echo "create group 'http-echo-user'"
  doPostRequest $jwtdAdminToken https://jwtd/project/default/group '{"name":"http-echo-user","rights":{"http-echo":{"role":"user"}}}'
  echo "get group 'http-echo-admin'"
  doGetRequest $jwtdAdminToken https://jwtd/project/default/group/http-echo-admin
  echo "add user 'admin' to 'http-echo-admin' and 'http-echo-user' group"
  doPatchRequest $jwtdAdminToken https://jwtd/project/default/user/admin '{"groups":["jwtd-admin","http-echo-admin","http-echo-user"]}'
  echo "add user 'example' to 'http-echo-user' group"
  doPatchRequest $jwtdAdminToken https://jwtd/project/default/user/example '{"groups":["http-echo-user"]}'
}

function testIt {
  echo "get role:admin token as admin"
  adminAdminToken=$(getToken admin admin http-echo role admin)
  echo "get role:admin token as user"
  userAdminToken=$(getToken example example http-echo role admin)
  echo "get role:user token as admin"
  adminUserToken=$(getToken admin admin http-echo role user)
  echo "get role:user token as user"
  userUserToken=$(getToken example example http-echo role user)

  echo "request http-echo /admin with admin-admin token"
  doGetRequest "${adminAdminToken}" https://http-echo/admin
  echo "request http-echo /admin with user-admin token"
  doGetRequest "${userAdminToken}" https://http-echo/admin
  echo "request http-echo /user with admin-user token"
  doGetRequest "${adminUserToken}" https://http-echo/user
  echo "request http-echo /user with user-user token"
  doGetRequest "${userUserToken}" https://http-echo/user
}

setupDocker 2>&1 > /dev/null
setupUserAndGroup
testIt
