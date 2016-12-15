#!/bin/bash

go build -ldflags '-linkmode external -extldflags -static' || exit $?
pushd jwtd-ctl >/dev/null
go build -ldflags '-linkmode external -extldflags -static' || exit $?
popd >/dev/null

docker build -t trusch/jwtd . || exit $?

exit $?
