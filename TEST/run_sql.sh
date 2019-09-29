#!/bin/sh
set -x
../../BASIL/basil --sql test.bs
../cyclone mysql 127.0.0.1 3306 TEST root root test.sql test_data.sql
