#!/bin/bash
. install.sh

echo "cornercheck must be checkouted in \"Go/src/github.com/sdolard/\" directory" 

echo "Run test..."

# go test -test.v . ./regions ./annonce # verbose version
go test . ./regions ./annonce
