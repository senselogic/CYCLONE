..\..\BASIL\basil --prefix postgresql_ --postgresql test.bs
..\cyclone --driver postgresql --host 127.0.0.1 --port 5432 --user root --password root postgresql_test.sql
pause
