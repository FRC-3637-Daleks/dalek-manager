#!/bin/bash
if [ ! -d cross-bin ]; then
    mkdir cross-bin
fi
export GOPATH=`pwd`
env GOOS="linux" GOWARCH="arm"  GOARM=7 go build -o $GOPATH/cross-bin/robot-manager $GOPATH/src/main.go