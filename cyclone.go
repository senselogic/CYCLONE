package main;

// -- IMPORTS

import ( "database/sql"
    "fmt"
    "strconv"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
    "github.com/gocql/gocql"
    "io/ioutil"
    "os"
    "strings"
    "time" );

// -- VARIABLES

var IsCqlDatabase, IsMysqlDatabase, IsPostgresqlDatabase, IsSqlDatabase bool;
var DatabaseDriver, DatabaseName, DatabasePassword, DatabasePort, DatabaseHost, DatabaseUser string;
var ExcludedCommandArray, ScriptFilePathArray [] string;
var CqlSession * gocql.Session;
var SqlDatabase * sql.DB;

// -- TYPES

type ERROR_MESSAGE struct {
    Text string;
    Error error;
}

// -- INQUIRIES

func ( error_message * ERROR_MESSAGE ) Print( ) {
    var error_text string;

    if ( error_message != nil ) {
        text := error_message.Text;

        if ( error_message.Error != nil ) {
            error_text = error_message.Error.Error();
        }

        if ( text != "" ) {
             if ( error_text != "" ) {
                fmt.Println( "*** ERROR :", text, "(" + error_text + ")" );
            } else {
                fmt.Println( "*** ERROR :", text );
            }
        } else if ( error_text != "" ) {
            fmt.Println( "*** ERROR :", error_text );
        }
    }
}

// -- OPERATIONS

func ( error_message * ERROR_MESSAGE ) Set( text string, error_ error ) {
    error_message.Text = text;
    error_message.Error = error_;
}

// ~~

func ( error_message * ERROR_MESSAGE ) SetText( text string ) {
    error_message.Text = text;
    error_message.Error = nil;
}

// ~~

func ( error_message * ERROR_MESSAGE ) SetError( error_ error ) {
    error_message.Text = "";
    error_message.Error = error_;
}

// -- FUNCTIONS

func IsNatural( text string ) bool {
    if ( len( text ) == 0 ) {
        return false;
    } else {
        for _, character := range text {
            if ( character < '0' || character > '9' ) {
                return false;
            }
        }

        return true;
    }
}

// ~~

func GetInteger( text string ) int {
    integer, _ := strconv.ParseInt( text, 10, 64 );

    return int( integer );
}

// ~~

func GetConditionalText( condition bool, text string ) string {
    if ( condition ) {
        return text;
    } else {
        return "" }
}

// ~~

func IsExcludedCommand( query string ) bool {
    for _, _excluded_command := range ExcludedCommandArray {
        if ( strings.HasPrefix( query, _excluded_command ) ) {
            return true;
        }
    }

    return false;
}

// ~~

func OpenDatabase( error_message * ERROR_MESSAGE ) bool {
    var error_ error;

    fmt.Println( "Opening database." );

    if ( IsCqlDatabase ) {
        cql_cluster_configuration := gocql.NewCluster( DatabaseHost );
        cql_cluster_configuration.Port = GetInteger( DatabasePort );
        cql_cluster_configuration.Timeout = 15 * time.Second;
        cql_cluster_configuration.ConnectTimeout = 15 * time.Second;
        cql_cluster_configuration.Consistency = gocql.Quorum;

        CqlSession, error_ = cql_cluster_configuration.CreateSession();
    } else if ( IsMysqlDatabase ) {
        SqlDatabase, error_ = sql.Open( "mysql", DatabaseUser + ":" + DatabasePassword + "@tcp(" + DatabaseHost + ":" + DatabasePort + ")/" + DatabaseName );
    } else if ( IsPostgresqlDatabase ) {
        SqlDatabase, error_ = sql.Open( "postgres", "host=" + DatabaseHost + " port=" + DatabasePort + " user=" + DatabaseUser + " password=" + DatabasePassword + GetConditionalText( DatabaseName != "", " dbname=" + DatabaseName ) + " sslmode=disable" );
    }

    if ( error_ != nil ) {
        error_message.SetError( error_ );

        return false;
    } else {
        return true;
    }
}

// ~~

func RunDatabaseQuery( query string, error_message * ERROR_MESSAGE ) bool {
    var error_ error;

    fmt.Println( query );

    if ( IsCqlDatabase ) {
        error_ = CqlSession.Query( query ).Exec();
    } else if ( IsSqlDatabase ) {
        _, error_ = SqlDatabase.Exec( query );
    }

    if ( error_ != nil ) {
        error_message.SetError( error_ );

        return false;
    } else {
        return true;
    }
}

// ~~

func ExecuteScripts( error_message * ERROR_MESSAGE ) bool {
    for _, script_file_path := range ScriptFilePathArray {
        fmt.Println( "Reading file : " + script_file_path )

        byte_array, error_ := ioutil.ReadFile( script_file_path );

        if ( error_ != nil ) {
            error_message.SetError( error_ );

            return false;
        }

        script := strings.ReplaceAll( string( byte_array ), "\r", "" );
        line_array := strings.Split( script, "\n" );
        query := "";

        for _, line := range line_array {
            if ( len( line ) > 0 && !strings.HasPrefix( line, "--" ) ) {
                line = strings.TrimSpace( line );

                if ( query == "" ) {
                    query = line;
                } else {
                    query += " " + line;
                }

                if ( strings.HasSuffix( query, ";" ) ) {
                    if ( !IsExcludedCommand( query ) && !RunDatabaseQuery( query, error_message ) ) {
                        return false;
                    }

                    query = "";
                }
            }
        }
    }

    return true;
}

// ~~

func CloseDatabase( ) bool {
    if ( IsCqlDatabase ) {
        CqlSession.Close();
    } else if ( IsSqlDatabase ) {
        SqlDatabase.Close();
    }

    return true;
}

// ~~

func ParseArguments( error_message * ERROR_MESSAGE ) bool {
    argument_array := os.Args[ 1 : ];

    for ( len( argument_array ) >= 1 && strings.HasPrefix( argument_array[ 0 ], "--" ) ) {
        if ( len( argument_array ) >= 2 && argument_array[ 0 ] == "--driver" ) {
            DatabaseDriver = argument_array[ 1 ];
            argument_array = argument_array[ 2 : ];
        } else if ( len( argument_array ) >= 2 && argument_array[ 0 ] == "--host" ) {
            DatabaseHost = argument_array[ 1 ];
            argument_array = argument_array[ 2 : ];
        } else if ( len( argument_array ) >= 2 && argument_array[ 0 ] == "--port" ) {
            DatabasePort = argument_array[ 1 ];
            argument_array = argument_array[ 2 : ];
        } else if ( len( argument_array ) >= 2 && argument_array[ 0 ] == "--user" ) {
            DatabaseUser = argument_array[ 1 ];
            argument_array = argument_array[ 2 : ];
        } else if ( len( argument_array ) >= 2 && argument_array[ 0 ] == "--password" ) {
            DatabasePassword = argument_array[ 1 ];
            argument_array = argument_array[ 2 : ];
        } else if ( len( argument_array ) >= 2 && argument_array[ 0 ] == "--database" ) {
            DatabaseName = argument_array[ 1 ];
            argument_array = argument_array[ 2 : ];
        } else if ( len( argument_array ) >= 2 && argument_array[ 0 ] == "--exclude" ) {
            ExcludedCommandArray = append( ExcludedCommandArray, argument_array[ 1 ] + " " );

            argument_array = argument_array[ 2 : ];
        } else {
            error_message.SetText( "Invalid option : " + argument_array[ 0 ] );

            return false;
        }
    }

    if ( len( argument_array ) >= 1 ) {
        ScriptFilePathArray = argument_array;

        fmt.Println( "Driver :", DatabaseDriver );
        fmt.Println( "Host :", DatabaseHost );
        fmt.Println( "Port :", DatabasePort );
        fmt.Println( "User :", DatabaseUser );
        fmt.Println( "Password :", DatabasePassword );
        fmt.Println( "Database :", DatabaseName );

        if ( DatabaseDriver == "cassandra" ) {
            IsCqlDatabase = true;
        } else if ( DatabaseDriver == "mysql" ) {
            IsSqlDatabase = true;
            IsMysqlDatabase = true;
        } else if ( DatabaseDriver == "postgresql" ) {
            IsSqlDatabase = true;
            IsPostgresqlDatabase = true;
        } else {
            error_message.SetText( "Invalid database driver : " + DatabaseDriver );

            return false;
        }

        if ( DatabaseHost == "" ) {
            error_message.SetText( "Invalid database server : " + DatabaseHost );

            return false;
        }

        if ( DatabasePort == "" || !IsNatural( DatabasePort ) ) {
            error_message.SetText( "Invalid database port : " + DatabasePort );

            return false;
        }

        if ( DatabaseUser == "" ) {
            error_message.SetText( "Missing database name argument : " + DatabaseUser );

            return false;
        }


        for _, script_file_path := range ScriptFilePathArray {
            if ( ( strings.HasSuffix( script_file_path, ".cql" ) && IsCqlDatabase ) || ( strings.HasSuffix( script_file_path, ".sql" ) && IsSqlDatabase ) ) {
                fmt.Println( "Script :", script_file_path );
            } else {
                error_message.SetText( "Invalid script argument : " + script_file_path );

                return false;
            }
        }
    } else {
        error_message.SetText( "Missing arguments" );

        return false;
    }

    return true;
}

// ~~

func main( ) {
    var error_message ERROR_MESSAGE;

    if ( ParseArguments( &error_message ) && OpenDatabase( &error_message ) && ExecuteScripts( &error_message ) && CloseDatabase() ) {
        fmt.Println( "Success." );
    } else {
        error_message.Print();
    }
}
