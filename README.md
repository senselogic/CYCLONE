![](https://github.com/senselogic/CYCLONE/blob/master/LOGO/cyclone.png)

# Cyclone

CQL and SQL script runner.

## Description

Cyclone runs the CQL/SQL statements of one or several scripts on a database.

## Installation

Install the Go compiler.

Build the executable with the following command lines :

```bash
go get -u github.com/go-sql-driver/mysql
go get -u github.com/gocql/gocql
go build cyclone.go
```

## Command line

```bash
cyclone {server} {port} {schema} {user} {password} {script file path} {script file path} ...
```

### Examples

```bash
cyclone 127.0.0.1 9042 TEST cassandra cassandra test.cql test_data.cql
```

Runs the CQL statements of `test.cql` and `test_data.cql` on `localhost`.

```bash
cyclone 127.0.0.1 3306 TEST root root test.sql test_data.sql
```

Runs the SQL statements of `test.sql` and `test_data.sql` on `localhost`.

## Version

1.0

## Author

Eric Pelzer (ecstatic.coder@gmail.com).

## License

This project is licensed under the GNU General Public License version 3.

See the [LICENSE.md](LICENSE.md) file for details.
