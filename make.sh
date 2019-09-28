#!/bin/sh
set -x
../GENERIS/generis --process ./ ./ --trim --join --go
go build cyclone.go
