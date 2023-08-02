#!/bin/sh
set -x
go mod download github.com/go-sql-driver/mysql
go mod download github.com/gocql/gocql
go mod download github.com/golang/snappy
go mod download github.com/hailocab/go-hostpool
go mod download github.com/lib/pq
go mod download gopkg.in/inf.v0
