#!/bin/sh
set -x
../GENERIS/generis --process ./ ./ --trim --join
go build cyclone.go
