#!/bin/bash

shellpath=`dirname $0`
cd $shellpath
pwd=`pwd`
export GOPATH="$pwd"
echo "$GOPATH"
go env
go build -o ./release/main src/main/main.go&&echo "build done"
