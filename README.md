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
cyclone {driver} {server} {port} {user} {password} {first_script} {second_script} ...
```

### Examples

```bash
cyclone cassandra 127.0.0.1 9042 cassandra cassandra test.cql test_data.cql
```

Runs the CQL statements of `test.cql` and `test_data.cql` on a local Cassandra database server.

```bash
cyclone mysql 127.0.0.1 3306 root root test.sql test_data.sql
```

Runs the SQL statements of `test.sql` and `test_data.sql` on a local MySQL database server.

## Version

1.0

## Author

Eric Pelzer (ecstatic.coder@gmail.com).

## License

This project is licensed under the GNU General Public License version 3.

See the [LICENSE.md](LICENSE.md) file for details.
