#!/bin/bash
. install.sh

echo "Run test..."

# go test -test.v . ./regions ./annonce # verbose version
go test . ./regions ./annonce
