#!/bin/bash

export GOPATH=/tmp/jwtd-vulcand-middleware
go get -u -v github.com/vulcand/vulcand/vbundle
go get -u -v github.com/trusch/jwtd

rm -r $GOPATH/src/github.com/vulcand/vulcand/vendor

mkdir -p $GOPATH/src/github.com/trusch/jwtd/vulcand
pushd $GOPATH/src/github.com/trusch/jwtd/vulcand
$GOPATH/bin/vbundle init --middleware github.com/trusch/jwtd/middleware
popd

go get -v github.com/trusch/jwtd/vulcand/...
go install github.com/trusch/jwtd/vulcand/...

cp $GOPATH/bin/vulcand $GOPATH/bin/vctl .

exit 0
