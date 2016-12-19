#!/bin/bash

go build -ldflags '-linkmode external -extldflags -static' || exit $?
pushd jwtd-ctl >/dev/null
go build -ldflags '-linkmode external -extldflags -static' || exit $?
popd >/dev/null

docker build -t trusch/jwtd . || exit $?

pushd jwtd-proxy >/dev/null
go build -ldflags '-linkmode external -extldflags -static' || exit $?
docker build -t trusch/jwtd-proxy . || exit $?
popd >/dev/null


exit $?
