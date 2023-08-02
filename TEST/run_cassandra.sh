#!/bin/sh
set -x
../../BASIL/basil --prefix cassandra_ --cql test.bs
../cyclone --driver cassandra --host 127.0.0.1 --port 9042 --user cassandra --password cassandra cassandra_test.cql
