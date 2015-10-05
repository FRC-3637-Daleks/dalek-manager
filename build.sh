#!/bin/bash
if [ -d bin ]; then
    mkdir bin
fi
export GOPATH=`pwd`
go build -o $GOPATH/bin/robot-manager $GOPATH/src/main.go