#!/bin/sh
set -x
../../BASIL/basil --cql test.bs
../cyclone 127.0.0.1 9042 TEST cassandra cassandra test.cql test_data.cql
